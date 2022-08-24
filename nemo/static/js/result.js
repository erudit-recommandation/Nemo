"use strict";

function copyToClipboard(text_Class) {
  const Toastify = window.Toastify;

  let text = "";
  const text_elements = document.getElementsByClassName(text_Class);
  for (let el of text_elements) {
    text += el.textContent + " ";
  }
  text = text.substring(0, text.length - 1);
  navigator.clipboard.writeText(text);
  Toastify({
    text: "text copier dans le presse papier",
    duration: 1000,
    newWindow: true,
    close: true,
    gravity: "bottom", // `top` or `bottom`
    position: "center", // `left`, `center` or `right`
    stopOnFocus: true, // Prevents dismissing of toast on hover
    onClick: function () {}, // Callback after click
  }).showToast();
}
