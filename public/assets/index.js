(() => {
  var __create = Object.create;
  var __defProp = Object.defineProperty;
  var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
  var __getOwnPropNames = Object.getOwnPropertyNames;
  var __getProtoOf = Object.getPrototypeOf;
  var __hasOwnProp = Object.prototype.hasOwnProperty;
  var __commonJS = (cb, mod) => function __require() {
    return mod || (0, cb[__getOwnPropNames(cb)[0]])((mod = { exports: {} }).exports, mod), mod.exports;
  };
  var __copyProps = (to, from, except, desc) => {
    if (from && typeof from === "object" || typeof from === "function") {
      for (let key of __getOwnPropNames(from))
        if (!__hasOwnProp.call(to, key) && key !== except)
          __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
    }
    return to;
  };
  var __toESM = (mod, isNodeMode, target) => (target = mod != null ? __create(__getProtoOf(mod)) : {}, __copyProps(
    // If the importer is in node compatibility mode or this is not an ESM
    // file that has been converted to a CommonJS file using a Babel-
    // compatible transform (i.e. "__esModule" has not been set), then set
    // "default" to the CommonJS "module.exports" for node compatibility.
    isNodeMode || !mod || !mod.__esModule ? __defProp(target, "default", { value: mod, enumerable: true }) : target,
    mod
  ));

  // node_modules/theme-change/index.js
  var require_theme_change = __commonJS({
    "node_modules/theme-change/index.js"(exports, module) {
      function themeToggle() {
        var toggleEl = document.querySelector("[data-toggle-theme]");
        var dataKey = toggleEl ? toggleEl.getAttribute("data-key") : null;
        (function(theme = localStorage.getItem(dataKey ? dataKey : "theme")) {
          if (localStorage.getItem(dataKey ? dataKey : "theme")) {
            document.documentElement.setAttribute("data-theme", theme);
            if (toggleEl) {
              [...document.querySelectorAll("[data-toggle-theme]")].forEach((el) => {
                el.classList.add(toggleEl.getAttribute("data-act-class"));
              });
            }
          }
        })();
        if (toggleEl) {
          [...document.querySelectorAll("[data-toggle-theme]")].forEach((el) => {
            el.addEventListener("click", function() {
              var themesList = el.getAttribute("data-toggle-theme");
              if (themesList) {
                var themesArray = themesList.split(",");
                if (document.documentElement.getAttribute("data-theme") == themesArray[0]) {
                  if (themesArray.length == 1) {
                    document.documentElement.removeAttribute("data-theme");
                    localStorage.removeItem(dataKey ? dataKey : "theme");
                  } else {
                    document.documentElement.setAttribute("data-theme", themesArray[1]);
                    localStorage.setItem(dataKey ? dataKey : "theme", themesArray[1]);
                  }
                } else {
                  document.documentElement.setAttribute("data-theme", themesArray[0]);
                  localStorage.setItem(dataKey ? dataKey : "theme", themesArray[0]);
                }
              }
              [...document.querySelectorAll("[data-toggle-theme]")].forEach((el2) => {
                el2.classList.toggle(this.getAttribute("data-act-class"));
              });
            });
          });
        }
      }
      function themeBtn() {
        var btnEl = document.querySelector("[data-set-theme='']");
        var dataKey = btnEl ? btnEl.getAttribute("data-key") : null;
        (function(theme = localStorage.getItem(dataKey ? dataKey : "theme")) {
          if (theme != void 0 && theme != "") {
            if (localStorage.getItem(dataKey ? dataKey : "theme") && localStorage.getItem(dataKey ? dataKey : "theme") != "") {
              document.documentElement.setAttribute("data-theme", theme);
              var btnEl2 = document.querySelector("[data-set-theme='" + theme.toString() + "']");
              if (btnEl2) {
                [...document.querySelectorAll("[data-set-theme]")].forEach((el) => {
                  el.classList.remove(el.getAttribute("data-act-class"));
                });
                if (btnEl2.getAttribute("data-act-class")) {
                  btnEl2.classList.add(btnEl2.getAttribute("data-act-class"));
                }
              }
            } else {
              var btnEl2 = document.querySelector("[data-set-theme='']");
              if (btnEl2.getAttribute("data-act-class")) {
                btnEl2.classList.add(btnEl2.getAttribute("data-act-class"));
              }
            }
          }
        })();
        [...document.querySelectorAll("[data-set-theme]")].forEach((el) => {
          el.addEventListener("click", function() {
            document.documentElement.setAttribute("data-theme", this.getAttribute("data-set-theme"));
            localStorage.setItem(dataKey ? dataKey : "theme", document.documentElement.getAttribute("data-theme"));
            [...document.querySelectorAll("[data-set-theme]")].forEach((el2) => {
              el2.classList.remove(el2.getAttribute("data-act-class"));
            });
            if (el.getAttribute("data-act-class")) {
              el.classList.add(el.getAttribute("data-act-class"));
            }
          });
        });
      }
      function themeSelect() {
        var selectEl = document.querySelector("select[data-choose-theme]");
        var dataKey = selectEl ? selectEl.getAttribute("data-key") : null;
        (function(theme = localStorage.getItem(dataKey ? dataKey : "theme")) {
          if (localStorage.getItem(dataKey ? dataKey : "theme")) {
            document.documentElement.setAttribute("data-theme", theme);
            var optionToggler = document.querySelector("select[data-choose-theme] [value='" + theme.toString() + "']");
            if (optionToggler) {
              [...document.querySelectorAll("select[data-choose-theme] [value='" + theme.toString() + "']")].forEach((el) => {
                el.selected = true;
              });
            }
          }
        })();
        if (selectEl) {
          [...document.querySelectorAll("select[data-choose-theme]")].forEach((el) => {
            el.addEventListener("change", function() {
              document.documentElement.setAttribute("data-theme", this.value);
              localStorage.setItem(dataKey ? dataKey : "theme", document.documentElement.getAttribute("data-theme"));
              [...document.querySelectorAll("select[data-choose-theme] [value='" + localStorage.getItem(dataKey ? dataKey : "theme") + "']")].forEach((el2) => {
                el2.selected = true;
              });
            });
          });
        }
      }
      function themeChange2(attach = true) {
        if (attach === true) {
          document.addEventListener("DOMContentLoaded", function(event) {
            themeToggle();
            themeSelect();
            themeBtn();
          });
        } else {
          themeToggle();
          themeSelect();
          themeBtn();
        }
      }
      if (typeof exports != "undefined") {
        module.exports = { themeChange: themeChange2 };
      } else {
        themeChange2();
      }
    }
  });

  // app/assets/index.js
  var import_theme_change = __toESM(require_theme_change());
  (0, import_theme_change.themeChange)();
  document.addEventListener("updateLeaderboard", function() {
    const data = JSON.parse(document.getElementById("leaderboardData").textContent);
    const names = data.map((row) => row.Name);
    const points = data.map((row) => row.Points);
    const ctx = document.getElementById("leaderboardChart").getContext("2d");
    new Chart(ctx, {
      type: "bar",
      data: {
        labels: names,
        datasets: [{
          label: "Xp per team",
          data: points,
          backgroundColor: "rgba(75, 192, 192, 0.2)",
          borderColor: "rgba(75, 192, 192, 1)",
          borderWidth: 1
        }]
      },
      options: {
        scales: {
          y: {
            beginAtZero: true
          }
        }
      }
    });
  });
})();
