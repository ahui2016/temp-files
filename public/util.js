const Util = {}

const DATE_TIME_FORMAT = 'YYYY-MM-DD HH:mm:ss';

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

function link(href) {
  return ce("a").attr({href: href});
}

/**
 * jElem is a string or JQuery<HTMLElement>
 */
function div(jElem) {
  if (typeof jElem == "string") {
    return ce("div").text(jElem);
  }
  return ce("div").append(jElem);
}

function timeNow() {
  return dayjs().format("HH:mm:ss");
}

Util.focus = function (jElem, timeout = 300) {
  setTimeout(() => {
    jElem.trigger("focus");
    jElem[0].setSelectionRange(1000, 1000);
  }, timeout);
};

/**
 * https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch
 * fetchOptions {url, obj, onSuccess, onError, onAlways}
 */
Util.fetch = async function (fetchOptions) {
  let resp = null;
  try {
    resp = await fetch(fetchOptions.url, fetchOptions.obj);
    if (!resp.ok) {
      let errMsg = `[${resp.status}] ${resp.statusText}`;
      if (resp.headers.get("Content-Type").startsWith("text/plain")) {
        const data = await resp.text();
        errMsg += ` ${data}`;
      }
      if (resp.headers.get("Content-Type").startsWith("application/json")) {
        const data = await resp.json(); 
        errMsg += ` ${JSON.stringify(data)}`;
      }
      throw new Error(errMsg);
    }
    if (fetchOptions.onSuccess) {
      fetchOptions.onSuccess(resp);
    }
  } catch (err) {
    if (fetchOptions.onError) {
      fetchOptions.onError(err);
    }
    console.error(err);
  }
  if (fetchOptions.onAlways) {
    fetchOptions.onAlways(resp);
  }
}

/**
 * fetchOptions {url, null, onSuccess, onError, onAlways}
 */
Util.postJSON = function (data, fetchOptions) {
  const obj = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  };
  fetchOptions.obj = obj;
  Util.fetch(fetchOptions);
}

/**
 * handlers {onSuccess, onError}
 */
Util.checkLogin = function (handlers) {
  Util.fetch({
    url: "/check-login",
    onSuccess: () => {
      if (handlers.onSuccess) handlers.onSuccess();
    },
    onError: (err) => {
      if (handlers.onError) handlers.onError(err);
    }
  });
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

/**
 * 获取地址栏的参数。
 * @param {string} name
 * @returns {string | null}
 */
function getUrlParam(name) {
  const queryString = new URLSearchParams(document.location.search);
  return queryString.get(name);
}
