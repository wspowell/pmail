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
    window.localStorage.setItem("jwt_token", jwtToken);
};
appStorage.prototype.JwtToken = function () {
    const jwtToken = window.localStorage.getItem("jwt_token");
    if (!jwtToken) { return null; }
    return jwtToken;
};
appStorage.prototype.UserGuid = function () {
    const jwtToken = this.JwtToken();
    if (!jwtToken) { return null; }
    const tokenData = parseJwt(jwtToken);
    return tokenData.user_guid;
};
appStorage.prototype.SetUser = function (userData) {
    const user = {
        user_guid: userData.user_guid,
        username: userData.username,
        mail_carry_capacity: userData.mail_carry_capacity,
        pineapple_on_pizza: userData.pineapple_on_pizza ? "Duh!" : "Party Pooper :/",
        mailbox_guid: userData.mailbox_guid,
    };

    window.localStorage.setItem("user_data", JSON.stringify(user));
};
appStorage.prototype.User = function () {
    const userData = window.localStorage.getItem("user_data");
    if (!userData) { return null; }
    return JSON.parse(userData);
};
appStorage.prototype.SetUserMailboxGuid = function (userMailboxGuid) {
    const user = this.User();
    if (!user) { return null; }
    user.mailbox_guid = userMailboxGuid;
    window.localStorage.setItem("user_data", JSON.stringify(user));
};
appStorage.prototype.SetUserMailbox = function (mailboxData) {
    const mailbox = {
        mailbox_guid: mailboxData.mailbox_guid,
        label: mailboxData.label,
        latitude: mailboxData.location.latitude,
        longitude: mailboxData.location.longitude,
        capacity: mailboxData.capacity,
        owner: mailboxData.owner,
    }
    window.localStorage.setItem("user_mailbox", JSON.stringify(mailbox));
};
appStorage.prototype.UserMailbox = function () {
    const mailbox = window.localStorage.getItem("user_mailbox");
    if (!mailbox) { return null; }
    return JSON.parse(mailbox);
};
appStorage.prototype.Logout = function () {
    window.localStorage.removeItem("jwt_token");
    window.localStorage.removeItem("user_data");
    window.localStorage.removeItem("user_mailbox");
};
const storage = new appStorage();