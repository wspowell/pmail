var SnailMail = {
    //baseUrl: "https://yrb7f0gokh.execute-api.us-east-1.amazonaws.com/test",
    baseUrl: "http://localhost:8080",

    sendRequest: function (method, path, requestBody, callbacks) {
        const Http = new XMLHttpRequest();
        const url = this.baseUrl + path;
        const body = JSON.stringify(requestBody);

        console.log("sending request -> " + method + " " + url + "\n" + body);

        let statusCode = 0;

        fetch(url, {
            method: method,
            headers: {
                "Content-Type": "application/json",
                "Accept": "application/json",
            },
            body: body
        }).then(response => {
            statusCode = response.status;
            if (callbacks[statusCode] == null) {
                throw new Error(method + " " + url + " missing callback for status code " + statusCode);
            }
            return response.json()
        }).then(body => {
            callbacks[statusCode](body)
        }).catch(error => {
            console.log(error.message)
        })
    },

    CreateUser: function (requestBody) {
        this.sendRequest("POST", "/users", requestBody, {
            201: function (body) {
                LoadContent("login", { info: "User successfully created! <happyface>" });
            },
            409: function (error) {
                console.log("user conflict: " + error.message);
                LoadContent("login", { info: "Username already exists <frustratedface>" });
            },
        });
    },

    AuthorizeUser: function (requestBody) {
        this.sendRequest("POST", "/authorize/user", requestBody, {
            200: function (body) {
                storage.SetJwtToken(body.jwt_token);
                LoadContent("user_info", storage.User());
            },
            404: function (error) {
                console.log("invalid login: " + error.message);
                LoadContent("login", { info: "Invalid username/password <sadface>" });
            },
        });
    }
};