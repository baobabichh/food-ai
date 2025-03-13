const good_nav_display = "";

function setCookie(name, value, hours) {
    var expires = "";
    if (hours) {
        var date = new Date();
        date.setTime(date.getTime() + (hours * 60 * 60 * 1000));
        expires = "; expires=" + date.toUTCString();
    }
    document.cookie = name + "=" + (value || "") + expires + "; path=/";
}
function getCookie(name) {
    var nameEQ = name + "=";
    var ca = document.cookie.split(';');
    for (var i = 0; i < ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) == ' ') c = c.substring(1, c.length);
        if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length, c.length);
    }
    return null;
}
function deleteCookie(name) {
    document.cookie = name + "=; Max-Age=0; path=/";
}

function setApiKey(value) {
    setCookie("apiKey", value, 1);
}

function getApiKey() {
    return getCookie("apiKey");
}

function deleteApiKey() {
    return deleteCookie("apiKey");
}

function updateNavLinks() {
    const loginLink = document.getElementById("nav_login");
    const logoutLink = document.getElementById("nav_logout");
    const registerLink = document.getElementById("nav_register");

    if (getApiKey()) {
        loginLink.style.display = 'none';
        registerLink.style.display = 'none';
        logoutLink.style.display = good_nav_display;
    } else {
        registerLink.style.display = good_nav_display;
        loginLink.style.display = good_nav_display;
        logoutLink.style.display = 'none';
    }
}

function onLogout() {
    deleteCookie("apiKey");
    updateNavLinks();
    window.location.href = "login.html";
}

document.getElementById("nav_logout").addEventListener('click', function (event) {
    event.preventDefault();
    onLogout();
});

updateNavLinks();

if (!getApiKey()) {
    const url = window.location.href;
    const urlParts = url.split('/');
    const fileName = urlParts[urlParts.length - 1] || 'index.html';

    if (fileName === "index.html" || fileName === "login.html" || fileName === "register.html") {

    } else {
        window.location.href = "index.html";
    }
}