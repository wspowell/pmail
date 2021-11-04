window.onload = function () {
    if (storage.LoggedIn()) {
        LoadContent("user_info", storage.User());
    } else {
        LoadContent("login");
    }
};

function Logout() {
    storage.Logout();
    LoadContent('login');
}

function getCreateUserData() {
    username = document.getElementById("username").value
    password = document.getElementById("password").value
    pineappleOnPizza = document.getElementById("pineapple_on_pizza").value
    if (pineappleOnPizza === "on") {
        pineappleOnPizza = true
    } else {
        pineappleOnPizza = false
    }

    return {
        username: username,
        password: password,
        pineapple_on_pizza: pineappleOnPizza
    }
}

function getAuthorizeUserData() {
    username = document.getElementById("username").value
    password = document.getElementById("password").value

    return {
        username: username,
        password: password,
    }
}