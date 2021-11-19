window.onload = function () {
    switch (window.location.hash) {
        case "#login":
            ShowLogin();
            break;
        case "#profile":
            if (storage.LoggedIn()) {
                ShowProfile();
            } else {
                ShowLogin();
            }
            break;
        case "#signup":
            LoadContent("create_user");
            break;
        default:
            if (storage.LoggedIn()) {
                ShowProfile();
            } else {
                ShowLogin();
            }
    }
};

function ShowLogin(pageData) {
    LoadContent("login", pageData);
}

function ShowCreateUser() {
    LoadContent("create_user");
}

function Login() {
    const loginInfo = {
        username: document.getElementById("username").value,
        password: document.getElementById("password").value,
    };

    SnailMail.AuthorizeUser(loginInfo, function (statusCode, responseBody) {
        switch (statusCode) {
            case 200:
                storage.SetJwtToken(responseBody.jwt_token);
                ShowProfile();
                break;
            case 404:
                UpdateContent({ error: "Invalid username/password <sadface>" });
                break;
            case 500:
            default:
                UpdateContent({ error: "Server slipped on a banana <surprisedface>" });
                break;
        }
    });
}

function Logout() {
    storage.Logout();
    LoadContent('login');
}

function CreateUser() {
    const userData = {
        username: document.getElementById("username").value,
        password: document.getElementById("password").value,
        pineapple_on_pizza: document.getElementById("pineapple_on_pizza").checked
    };

    SnailMail.CreateUser(userData, function (statusCode, responseBody) {
        switch (statusCode) {
            case 201:
                LoadContent("login", { info: "User successfully created! <happyface>" });
                break;
            case 409:
                UpdateContent({ error: "Username already exists <frustratedface>" });
                break;
            case 500:
            default:
                UpdateContent({ error: "Server slipped on a banana <surprisedface>" });
                break;
        }
    });
}

function ShowProfile() {
    let userGuid = storage.UserGuid();
    let user = storage.User();
    let userMailbox = storage.UserMailbox();

    if (user) {
        showUser();
    } else {
        if (userGuid) {
            SnailMail.GetUser(userGuid, function (statusCode, responseBody) {
                switch (statusCode) {
                    case 200:
                        user = responseBody;
                        user.user_guid = userGuid;
                        storage.SetUser(responseBody);
                        user = storage.User();
                        showUser();
                        break;
                    case 404:
                        userMailboxElement.classList.add('hidden');
                        createUserMailboxElement.classList.remove("hidden");
                        //UpdateContent({ error: "Mailbox not found <confusedface>" });
                        break;
                    case 500:
                    default:
                        UpdateContent({ error: "Server slipped on a banana <surprisedface>" });
                }
            });
        } else {
            ShowLogin({ error: "Must be logged in <sternface>" });
        }
    }

    function showUser() {
        LoadContent("user_info", user);

        const userMailboxElement = document.getElementById("user_mailbox");
        const createUserMailboxElement = document.getElementById("create_user_mailbox_button");

        if (user.mailbox_guid) {
            userMailboxElement.classList.remove("hidden");
            createUserMailboxElement.classList.add("hidden");

            showMailbox();
        } else {
            userMailboxElement.classList.add('hidden');
            createUserMailboxElement.classList.remove("hidden");
        }
    }

    function showMailbox() {

        if (userMailbox) {
            UpdateContent(userMailbox);
        } else {
            SnailMail.GetMailbox(user.mailbox_guid, function (statusCode, responseBody) {
                switch (statusCode) {
                    case 200:
                        storage.SetUserMailbox(responseBody);
                        userMailbox = storage.UserMailbox();
                        UpdateContent(userMailbox);
                        break;
                    case 404:
                        UpdateContent({ error: "Mailbox not found <confusedface>" });
                        break;
                    case 500:
                    default:
                        UpdateContent({ error: "Server slipped on a banana <surprisedface>" });
                }
            });
        }
    }
}

function ShowCreateUserMailbox() {
    LoadContent("create_user_mailbox");
}

function CreateUserMailbox() {
    const userGuid = storage.UserGuid();

    if (!userGuid) {
        ShowLogin({ error: "Must be logged in <sternface>" });
        return;
    }

    const user = storage.User();

    const mailboxData = {
        label: document.getElementById("label").value,
        location: {
            latitude: parseFloat(document.getElementById("latitude").value),
            longitude: parseFloat(document.getElementById("longitude").value),
        },
        capacity: 20,//document.getElementById("capacity").value,
        owner: user.user_guid,
    };

    SnailMail.CreateMailbox(mailboxData, function (statusCode, responseBody) {
        switch (statusCode) {
            case 201:
                storage.SetUserMailboxGuid(responseBody.mailbox_guid);
                mailboxData.mailbox_guid = responseBody.mailbox_guid;
                storage.SetUserMailbox(mailboxData);
                ShowProfile();
                break;
            case 404:
                UpdateContent({ error: "User not found <confusedface>" });
                break;
            case 409:
                UpdateContent({ error: "Mailbox already exists <confusedface>" });
                break;
            case 500:
            default:
                UpdateContent({ error: "Server slipped on a banana <surprisedface>" });
        }
    });
}