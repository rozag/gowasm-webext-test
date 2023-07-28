function headersArrayToObject(
  headersArray: chrome.webRequest.HttpHeader[]
): Record<string, string> {
  const headersObject: Record<string, string> = {};

  for (const header of headersArray) {
    if (header.name && header.value) {
      headersObject[header.name] = header.value;
    }
  }

  return headersObject;
}

const _initiatorToUrlPatternsMap: Map<string, Set<string>> = new Map([
  ["https://app.wingman.wtf", new Set(["/hot-flights"])],
]);

chrome.webRequest.onBeforeSendHeaders.addListener(
  function(
    details: chrome.webRequest.WebRequestHeadersDetails
  ): chrome.webRequest.BlockingResponse | void {
    console.log("beforeSendHeaders: details:", details);

    if (!details.initiator) {
      console.log("beforeSendHeaders: abort: no initiator");
      return;
    }

    const patterns = _initiatorToUrlPatternsMap.get(details.initiator);
    if (!patterns) {
      console.log("beforeSendHeaders: abort: not our initiator");
      return;
    }

    let isMatchingUrl = false;
    for (const pattern of patterns) {
      if (details.url.includes(pattern)) {
        isMatchingUrl = true;
        break;
      }
    }
    if (!isMatchingUrl) {
      console.log("beforeSendHeaders: abort: not our URL pattern");
      return;
    }

    if (details.method !== "GET") {
      // TODO: other methods should be supported later
      console.log("beforeSendHeaders: abort: not GET method");
      return;
    }

    let headers = {};
    if (details.requestHeaders) {
      headers = headersArrayToObject(details.requestHeaders);
    }
    console.log("beforeSendHeaders: headers:", headers);

    fetch(details.url, {
      method: details.method,
      headers: headers,
    })
      .then((response) => response.text())
      .then((text) => console.log("beforeSendHeaders: response text:", text));
  },
  { urls: ["<all_urls>"] },
  ["requestHeaders", "extraHeaders"]
);
