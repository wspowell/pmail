
// Set your needed values
var IDENTITY_POOL_ID = 'cognito-idp.us-east-1.us-east-1.amazoncognito.com/us-east-1_jsehCKXod';
var ACCOUNT_ID = '677645055400';
var REGION = 'us-east-1';

// Initialize the Amazon Cognito credentials provider
AWS.config.region = REGION; // Region
AWS.config.credentials = new AWS.CognitoIdentityCredentials({ IdentityPoolId: IDENTITY_POOL_ID, });

var getIdParams = {
    IdentityPoolId: IDENTITY_POOL_ID,
    AccountId: ACCOUNT_ID
};

var cognitoidentity = new AWS.CognitoIdentity({ apiVersion: '2014-06-30' });

cognitoidentity.getId(getIdParams, function (err, data) {
    if (err) {
        results.innerHTML = "Error " + err;
    } else {
        results.innerHTML = "Cognito Identity ID is " + data.IdentityId;
    }
});