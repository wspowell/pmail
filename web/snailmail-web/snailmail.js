var SnailMail = {
    //baseUrl: "https://yrb7f0gokh.execute-api.us-east-1.amazonaws.com/test",
    baseUrl: "http://localhost:8080",

    sendRequest: function (method, path, requestBody, callback) {
        const Http = new XMLHttpRequest();
        const url = this.baseUrl + path;
        const body = JSON.stringify(requestBody);

        console.debug("sending request -> " + method + " " + url + "\n" + (body ? body : "{}"));

        let statusCode = 0;

        const requestData = {
            method: method,
            headers: {
                "Content-Type": "application/json",
                "Accept": "application/json",
            },
        }

        if (requestBody) {
            requestData.body = body;
        }

        fetch(url, requestData).then(response => {
            statusCode = response.status;
            return response.json()
        }).then(responseBody => {
            console.debug("response " + statusCode + " " + JSON.stringify(responseBody));
            callback(statusCode, responseBody)
        }).catch(error => {
            console.log(error.message)
            callback(500, error)
        }).finally(() => {
            //
        })
    },

    CreateUser: function (requestBody, callback) {
        this.sendRequest("POST", "/users", requestBody, callback)
    },

    GetUser: function (userGuid, callback) {
        this.sendRequest("GET", "/users/" + userGuid, null, callback)
    },

    AuthorizeUser: function (requestBody, callback) {
        this.sendRequest("POST", "/authorize/user", requestBody, callback)
    },

    CreateMailbox: function (requestBody, callback) {
        this.sendRequest("POST", "/mailboxes", requestBody, callback);
    },

    GetMailbox: function (mailboxGuid, callback) {
        this.sendRequest("GET", "/mailboxes/" + mailboxGuid, null, callback)
    },
};