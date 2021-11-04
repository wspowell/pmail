function parseJwt(token) {
    var base64Url = token.split('.')[1];
    var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    var jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
};

function appStorage() { };
appStorage.prototype.LoggedIn = function () {
    if (window.localStorage.getItem("jwt_token") != null) {
        return true;
    }
    return false;
};
appStorage.prototype.SetJwtToken = function (jwtToken) {
    return window.localStorage.setItem("jwt_token", jwtToken);
};
appStorage.prototype.User = function () {
    const jwtToken = window.localStorage.getItem("jwt_token");

    if (jwtToken == null) {
        return null;
    }

    const data = parseJwt(jwtToken);

    return {
        UserGuid: data.user_guid,
        Username: data.username,
        MailCarryCapacity: data.mail_carry_capacity,
        PineappleOnPizza: data.pineapple_on_pizza,
    };
};
appStorage.prototype.Logout = function () {
    window.localStorage.removeItem("jwt_token");
};
let storage = new appStorage();