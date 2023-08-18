const Util = {}

/**
 * jElem is a JQuery<HTMLElement>
 */
Util.disable = function (jElem) {
  const nodeName = jElem.prop("nodeName");
  if (nodeName == "BUTTON" || nodeName == "INPUT") {
    jElem.prop("disabled", true);
  } else {
    jElem.css("pointer-events", "none");
  }
};

/**
 * jElem is a JQuery<HTMLElement>
 */
Util.enable = function (jElem) {
  const nodeName = jElem.prop("nodeName");
  if (nodeName == "BUTTON" || nodeName == "INPUT") {
    jElem.prop("disabled", false);
  } else {
    jElem.css("pointer-events", "auto");
  }
};

/**
 * Get a raw element using document.querySelector
 */
function raw(selectors) {
  return document.querySelector(selectors);
}

/**
 * ce: create element
 */
function ce(tagName) {
  return $(document.createElement(tagName));
}

function br() {
  return ce("br");
}

function span(text) {
  if (typeof text == "string") {
    return ce("span").text(text);
  }
  return ce("span");
}

function p(text) {
  if (typeof text == "string") {
    return ce("p").text(text);
  }
  return ce("p");
}

function timeNow() {
  return dayjs().format("HH:mm:ss");
}

/**
 * Upload file using an HTML input element, FormData and fetch().
 * https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch
 * fetchOptions {url, formData, onSuccess, onError, onAlways}
 */
async function uploadFile(fetchOptions) {
  let resp = null;
  try {
    resp = await fetch(fetchOptions.url, {
      method: "POST",
      body: fetchOptions.formData,
    });
    if (!resp.ok) {
      throw new Error(`Response NG: [${resp.status}] ${resp.statusText}`);
    }
    if (fetchOptions.onSuccess) {
      fetchOptions.onSuccess(resp);
    }
  } catch (err) {
    if (fetchOptions.onError) {
      fetchOptions.onError(err);
    }
    console.error("Error:", err);
  }
  if (fetchOptions.onAlways) {
    fetchOptions(resp);
  }
}

/**
 * 把文檔體積轉換為方便人類閱讀的形式。
 */
function fileSizeToString(fileSize, fixed) {
  if (fixed == null) {
    fixed = 2;
  }
  const sizeGB = fileSize / 1024 / 1024 / 1024;
  if (sizeGB < 1) {
    const sizeMB = sizeGB * 1024;
    if (sizeMB < 1) {
      const sizeKB = sizeMB * 1024;
      return `${sizeKB.toFixed(fixed)} KB`;
    }
    return `${sizeMB.toFixed(fixed)} MB`;
  }
  return `${sizeGB.toFixed(fixed)} GB`;
}
