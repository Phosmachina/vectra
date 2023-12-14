// noinspection DuplicatedCode

//region CONSTANTS

const TOKEN_LENGTH = 45

const NOTIFICATION_TIMEOUT = 5000

const NOTIFICATION_TYPE = Object.freeze({
    INFO: "notification-info",
    WARN: "notification-warn",
    ERROR: "notification-error",
})

//endregion

document.addEventListener("DOMContentLoaded", () => {
    surveyForm()
    loadSvgSprite()
    loadSystemTheme()
});

function loadSvgSprite() {

    let loc = window.location

    fetch(`${loc.protocol}//static.${loc.hostname}:${loc.port}/svg/sprite`)
        .then(response => response.text())
        .then(data => {
            document.querySelector("#svg-sprite").outerHTML = data
        })
}

function loadSystemTheme() {
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: light)').matches) {
        toggleTheme()
    }

    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', event => {
        const newColorScheme = event.matches ? "dark" : "light"
        document.querySelector("html").className = `theme-${newColorScheme}`

    })
}

const handleInput = (labels, idxLabel) => {
    let input = labels[idxLabel].querySelector("input") || labels[idxLabel].querySelector("textarea")
    if (!input) return null

    let allVerify = Array(5)

    for (let i = 0; i < 5; i++) {
        let funcName = input.getAttribute(`validation-${i}`)
        if (!funcName) break

        let verify = getVerifyFunction(labels, idxLabel, input, funcName, i)
        allVerify[i] = verify

        if (funcName === "sameAsPrevious" && idxLabel > 0) {
            let prev = labels[idxLabel - 1].querySelector("input")
            let prevOnChange = prev.onchange
            prev.onchange = () => verify() & prevOnChange()
        }
    }

    input.onchange = () => allVerify.forEach(verify => verify && verify())

    return labels[idxLabel]
}

const surveyForm = () => {
    let forms = document.querySelectorAll("form")

    for (let form of forms) {

        let labels = form.querySelectorAll("label")
        let inputs = Array.from(
            {length: labels.length},
            (_, idxLabel) => handleInput(labels, idxLabel)
        ).filter(input => input !== null)

        form.onsubmit = (event) => {
            if (event) event.preventDefault()

            let isReady = inputs.every(label => {
                label.querySelector("input").onchange(null)
                return label.querySelectorAll(".advice").length === 0
            })

            // for (let file of form.querySelector(".hide_file")) {
            //     if (!window[file.getAttribute("validation-func")](file)) isReady = false
            // }

            if (isReady) {
                let split = form.getAttribute("ajax-func").split(" ")
                window[split[0]](...split.slice(1))
            }
        }
    }
}

function getVerifyFunction(labels, idxLabel, input, funcName, index) {

    return function () {
        let isValid = funcName === "sameAsPrevious" ?
            window[funcName](labels[idxLabel - 1]
                .querySelector("input").value, input.value)
            : window[funcName](input.value)
        let advice_group = labels[idxLabel].querySelector(".advice-group")

        input.classList.remove(isValid ? "is-error" : "is-valid")
        input.classList.add(isValid ? "is-valid" : "is-error")

        let adviceNode = advice_group.querySelector(`.advice-${index}`)

        if (isValid) {
            if (adviceNode)
                advice_group.removeChild(adviceNode)
        } else if (!adviceNode) {
            let advice = document.querySelector("#advice").content.cloneNode(true).querySelector("div")
            advice.classList.add(`advice-${index}`)
            advice.innerHTML = input.getAttribute(`advice-${index}`)
            advice_group.appendChild(advice)
        }

        return isValid
    }
}

//region VALIDATION

/**
 * Checks if the given input element has any selected files.
 *
 * @param {HTMLElement} input - The input element to check.
 * @returns {boolean} - Returns true if the input element has files selected, otherwise false.
 */
function hasFile(input) {
    return input.files.length > 0
}

/**
 * Checks if the given input string is not empty and has a length of at least 3 characters.
 *
 * @param {string} input - The input string to be validated.
 * @return {boolean} - Returns true if the input string is not empty and has a length of at least 3 characters,
 *                     otherwise returns false.
 */
function notEmpty(input) {
    return input.trim().length >= 3
}

/**
 * Validates an email address.
 *
 * @param {string} email - The email address to validate.
 * @return {boolean} - True if the email address is valid, false otherwise.
 */
function validateEmail(email) {
    const res = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
    return res.test(String(email).toLowerCase().trim())
}

/**
 * Validates a phone number using a regular expression.
 *
 * @function validatePhone
 * @param {string|number} phone - The phone number to be validated.
 * @returns {boolean} - True if the phone number is valid, otherwise false.
 */
function validatePhone(phone) {
    const res = /^(?:(?:\+|00)33|0)\s*[1-9](?:[\s.-]*\d{2}){4}$/
    return res.test(String(phone).trim())
}

/**
 * Validates a zip code.
 *
 * @param {string} zip - The zip code to validate.
 * @return {boolean} - Returns true if the zip code is valid, false otherwise.
 */
function validateZip(zip) {
    const res = /^\d{5}$|^\d{5}-\d{4}$/
    return res.test(String(zip).trim())
}

/**
 * Determines whether the current value is the same as the previous value.
 *
 * @param {any} prev - The previous value.
 * @param {any} current - The current value.
 * @returns {boolean} - True if the current value is the same as the previous value; otherwise, false.
 */
function sameAsPrevious(prev, current) {
    return prev === current
}

/**
 * Validates the given initialization token based on its length.
 *
 * @param {string} token - The initialization token to be validated.
 * @note The token length of TOKEN_LENGTH characters doesn't guarantee its validity, further checks may be essential.
 * @return {boolean} - Returns true if the length of the token is as per the defined length criteria, false otherwise.
 */
function validateInitToken(token) {
    return token.trim().length === TOKEN_LENGTH
}

/**
 * Checks if the input value contains at least one uppercase letter.
 *
 * @param {string} value - The string to check.
 * @returns {boolean} Returns true if the value contains at least one uppercase letter, otherwise it returns false.
 */
function atLeastUpper(value) {
    return /[A-Z]/.test(value)
}

/**
 * Checks if the input value contains at least one lowercase letter.
 *
 * @param {string} value - The string to check.
 * @returns {boolean} Returns true if the value contains at least one lowercase letter, otherwise it returns false.
 */
function atLeastLower(value) {
    return /[a-z]/.test(value)
}

/**
 * Checks if the input value contains at least one digit.
 *
 * @param {string} value - The string to check.
 * @returns {boolean} Returns true if the value contains at least one digit, otherwise it returns false.
 */
function atLeastNumber(value) {
    return /\d/.test(value)
}

/**
 * Checks if the input value contains at least one special character (non-alphanumeric).
 *
 * @param {string} value - The string to check.
 * @returns {boolean} Returns true if the value contains at least one non-alphanumeric character, otherwise it returns false.
 */
function atLeastSpecialChar(value) {
    return /[^a-zA-Z0-9]/.test(value)
}

/**
 * Checks if the input value has at least 8 characters.
 *
 * @param {string} value - The string to check.
 * @returns {boolean} Returns true if the value has at least 8 characters, otherwise it returns false.
 */
function atLeast8Char(value) {
    return value.length >= 8
}

//endregion

//region FETCH

function initRequest() {
    simpleFetch("/api/v1/activate/admin", mapFormValues(document.querySelector("form")))
        .then(() => window.location = "/login")
        .catch(() => document.querySelector("[name='token']").value = "")
}

function login() {
    simpleFetch("/api/v1/login", mapFormValues(document.querySelector("form")))
        .then(() => window.location = "/")
        .catch(() => document.querySelector("[name='password']").value = "")
}

function adminLogin() {
    let map = mapFormValues(document.querySelector("form"))
    map.set("is_admin", true)

    simpleFetch("/api/v1/login", map)
        .then(() => window.location = "/admin")
        .catch(() => document.querySelector("[name='password']").value = "")
}

//endregion

//region BUTTONS

function toggleTheme() {
    let htmlTag = document.querySelector("html")
    htmlTag.className = htmlTag.classList.contains("theme-dark") ? "theme-light" : "theme-dark"

    document.querySelectorAll("#theme-switcher div")
        .forEach(blk => blk.classList.toggle("hidden"))
}

function toggleLang(btn) {
    let lang = btn.querySelector("p").innerText;
    simpleFetch("/api/v1/update/lang", new Map([["lang", lang]]))
        .then(() => document.location.reload())
        .catch()
}

//endregion

//region HELPERS

/**
 * Makes an asynchronous HTTP POST request to the specified URL with a content payload and CSRF token.
 *
 * @param {string} url - The URL to send the request to.
 * @param {Object} content - The payload of the request.
 * @returns {Promise} - A promise that resolves to the parsed JSON response from the server.
 * @throws {Error} - If the HTTP response status is not OK.
 */
async function fetchRequest(url, content) {
    const csrfToken = document.cookie
        .split('; ')
        .find(row => row.startsWith('csrf-token='))
        .split('=')[1]

    return fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'charset': 'UTF-8',
            'X-CSRF-Token': csrfToken
        },
        body: JSON.stringify(content)
    }).then(response => {
        return response.json()
    })
}

/**
 * Send a fetch request with the specified URL and payload.
 *
 * @param {string} url - The URL to send the fetch request to.
 * @param {Map} map - The map object containing the payload data.
 * @returns {Promise} A promise that resolves to the response of the fetch request. If the response has a reason property, an error is thrown with the reason.
 */
function simpleFetch(url, map) {
    return fetchRequest(url, Object.fromEntries(map))
        .then(response => {
            if (response.reason) {
                newNotification(response.reason, NOTIFICATION_TYPE.ERROR)
                throw new Error(`Server reject the request with reason: ${response.reason}`)
            }
            return response
        })
}

/**
 * Retrieves data from a given HTML form and returns it as a Map.
 *
 * @param {HTMLFormElement} form - The HTML form element from which to retrieve data.
 * @return {Map} - A Map with the form data, where the keys are the names of the form controls and the values are the corresponding values entered by the user.
 */
function mapFormValues(form) {

    let map = new Map()
    let controls = form.querySelectorAll("input, select, textarea")

    controls.forEach(control => {
        if (control.name.length === 0) return
        map.set(control.name, control.value)
    })

    return map
}

function newNotification(msg, type) {

    let notification_group = document.querySelector("#notification-group")
    let notification = document.querySelector("#notification").content.cloneNode(true).querySelector("div")
    let spans = notification.querySelectorAll("span")

    let closed = false
    let close = function () {
        if (!closed) {
            notification_group.removeChild(notification)
            closed = true
        }
    }

    notification.classList.add(type)
    spans[0].innerHTML = msg
    spans[1].onclick = close

    notification_group.appendChild(notification)

    setTimeout(close, NOTIFICATION_TIMEOUT)
}

//endregion
