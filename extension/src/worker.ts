import { initiatorToUrlPatternsMap } from "./targets.js";
import setUpWasm from "./wasm_exec.js";

setUpWasm();

const logOn = false;

function log(...args: any[]) {
  if (logOn) {
    console.log(...args);
  }
}

function runWasm() {
  // @ts-ignore
  const go = new Go();
  // const wasmPath = chrome.runtime.getURL("client-js-wasm.wasm");
  WebAssembly.instantiateStreaming(
    fetch("./client-js-wasm.wasm"),
    go.importObject
  ).then((result) => {
    go.run(result.instance);
  });
}

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

chrome.webRequest.onBeforeSendHeaders.addListener(
  function(
    details: chrome.webRequest.WebRequestHeadersDetails
  ): chrome.webRequest.BlockingResponse | void {
    log("beforeSendHeaders: details:", details);

    if (!details.initiator) {
      log("beforeSendHeaders: abort: no initiator");
      return;
    }

    const patterns = initiatorToUrlPatternsMap.get(details.initiator);
    if (!patterns) {
      log("beforeSendHeaders: abort: not our initiator");
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
      log("beforeSendHeaders: abort: not our URL pattern");
      return;
    }

    if (details.method !== "GET") {
      // TODO: other methods should be supported later
      log("beforeSendHeaders: abort: not GET method");
      return;
    }

    let headers = {};
    if (details.requestHeaders) {
      headers = headersArrayToObject(details.requestHeaders);
    }
    log("beforeSendHeaders: headers:", headers);

    fetch(details.url, {
      method: details.method,
      headers: headers,
    })
      .then((response) => response.text())
      .then((text) => log("beforeSendHeaders: response text:", text));

    runWasm();
  },
  { urls: ["<all_urls>"] },
  ["requestHeaders", "extraHeaders"]
);
