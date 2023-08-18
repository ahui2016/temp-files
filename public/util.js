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
