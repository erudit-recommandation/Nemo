"use strict";

const Toastify = window.Toastify;

class InteractiveInput {
  constructor() {
    const corpusElements = document.querySelectorAll("option");
    const selectedCorpusValue =
      document.getElementsByTagName("select")[0].value;
    this.selectedCorpus = {
      value: selectedCorpusValue,
      name: document.querySelector(`#${selectedCorpusValue}`).innerHTML,
    };

    this.indexCorpus = 0;

    this.corpusElement = document.querySelector(".corpus");
    this.corpusElement.textContent = this.selectedCorpus.name;

    this.corpus = Array.from(corpusElements).map((el) => {
      return {
        value: el.value,
        name: el.innerText,
      };
    });

    this.btnRencontreEnvoyage = document.querySelector(
      ".rencontre-en-voyage, .rencontre-en-voyage-previous-search"
    );
    this.btnEntenduEnVoyage = document.querySelector(
      ".entendu-en-voyage, .entendu-en-voyage-previous-search"
    );
    this.btnClefCanonique = document.querySelector(
      ".clef-canonique, .clef-canonique-previous-search"
    );
    this.searchBar = document.querySelector(".search-bar");

    document.querySelector(".input-form").onsubmit = () => {
      return false;
    };

    this.btnSelected = this.btnEntenduEnVoyage;
    this.deleteSelectTag();

    this.startingFunction();

    this.setOnclickEvent();
  }

  setOnclickEvent() {
    this.btnRencontreEnvoyage.onclick = () => {
      this.btnSelected = this.btnRencontreEnvoyage;
      this.rencontreEnVoyageInteraction();
    };
    this.btnEntenduEnVoyage.onclick = () => {
      this.btnSelected = this.btnEntenduEnVoyage;
      this.entenduEnVoyageInteraction();
    };

    this.btnClefCanonique.onclick = () => {
      this.btnSelected = this.btnClefCanonique;
      this.clefCanoniqueInteraction();
    };

    this.corpusElement.onclick = () => {
      this.rotateCorpus();
    };
  }

  search() {
    this.searchBar.addEventListener("keypress", (event) => {
      if (event.keyCode == 13 && !event.shiftKey) {
        event.preventDefault();
        if (this.searchBar.value == "") {
          const msg = "un texte doit être entré";
          try {
            Toastify({
              text: msg,
              duration: 1000,
              newWindow: true,
              close: true,
              gravity: "top", // `top` or `bottom`
              position: "center", // `left`, `center` or `right`
              stopOnFocus: true, // Prevents dismissing of toast on hover
              onClick: function () {}, // Callback after click
            }).showToast();
          } catch (err) {
            alert(msg);
          }
        } else {
          const form = document.createElement("form");
          form.style.visibility = "hidden";
          document.body.appendChild(form);
          const text = document.createElement("input");
          text.type = "hidden";
          text.name = "text";
          text.value = this.searchBar.value;

          const corpus = document.createElement("input");
          corpus.type = "hidden";
          corpus.name = "corpus";
          corpus.value = this.selectedCorpus.value;

          form.appendChild(text);
          form.appendChild(corpus);
          form.method = "POST";
          form.action = this.btnSelected.formAction;
          form.submit();
        }
      }
    });
  }

  deleteSelectTag() {
    document.querySelector("select").style.visibility = "hidden";
    document.querySelector("select").style.marginTop = "0px";
    document.querySelector("select").style.marginBottom = "0px";
  }

  rencontreEnVoyageInteraction() {
    this.setOneButtonBlack(this.btnRencontreEnvoyage);
    if (this.searchBar.tagName === "TEXTAREA") {
      this.searchBar.style.height = "140px";
    } else {
      const d = document.createElement("TEXTAREA");
      d.value = this.searchBar.value;
      d.classList = this.searchBar.classList;
      this.searchBar.parentNode.replaceChild(d, this.searchBar);

      this.searchBar = document.querySelector(".search-bar");
      this.searchBar.style.height = "140px";
      this.search();
    }
  }

  entenduEnVoyageInteraction() {
    this.setOneButtonBlack(this.btnEntenduEnVoyage);
    if (this.searchBar.tagName === "TEXTAREA") {
      const d = document.createElement("INPUT");
      d.value = this.searchBar.value;
      d.classList = this.searchBar.classList;

      this.searchBar.parentNode.replaceChild(d, this.searchBar);
      this.searchBar = document.querySelector(".search-bar");
      this.search();
    }
  }

  clefCanoniqueInteraction() {
    this.setOneButtonBlack(this.btnClefCanonique);
    const msg = "Service à venir!";
    try {
      Toastify({
        text: msg,
        duration: 1000,
        newWindow: true,
        close: true,
        gravity: "top", // `top` or `bottom`
        position: "center", // `left`, `center` or `right`
        stopOnFocus: true, // Prevents dismissing of toast on hover
        onClick: function () {}, // Callback after click
      }).showToast();
    } catch (err) {
      alert(msg);
    }
  }

  rotateCorpus() {
    this.indexCorpus += 1;
    const i = this.indexCorpus % this.corpus.length;
    this.selectedCorpus = this.corpus[i];
    this.corpusElement.textContent = this.selectedCorpus.name;
    console.log(this.selectedCorpus);
  }

  setOneButtonBlack(button) {
    if (button === this.btnRencontreEnvoyage) {
      this.btnRencontreEnvoyage.style.backgroundImage =
        'url("/static/images/o_black.svg")';
    } else {
      this.btnRencontreEnvoyage.style.backgroundImage =
        'url("/static/images/o_white.svg")';
    }

    if (button === this.btnEntenduEnVoyage) {
      this.btnEntenduEnVoyage.style.backgroundImage =
        'url("/static/images/p_black.svg")';
    } else {
      this.btnEntenduEnVoyage.style.backgroundImage =
        'url("/static/images/p_white.svg")';
    }

    if (button === this.btnClefCanonique) {
      this.btnClefCanonique.style.backgroundImage =
        'url("/static/images/key_black.svg")';
    } else {
      this.btnClefCanonique.style.backgroundImage =
        'url("/static/images/key_white.svg")';
    }
  }

  startingFunction() {
    const url = window.location.href;
    if (url.includes(this.btnClefCanonique.formAction)) {
      this.clefCanoniqueInteraction();
    } else if (url.includes(this.btnEntenduEnVoyage.formAction)) {
      this.entenduEnVoyageInteraction();
    } else if (url.includes(this.btnRencontreEnvoyage.formAction)) {
      this.rencontreEnVoyageInteraction();
    } else {
      this.entenduEnVoyageInteraction();
    }
  }
}

const interactiveInput = new InteractiveInput();
