var app = angular.module('main', ['angular-loading-bar']);

$(function() {
    $("#url").focus();
    $("#encryption").focus();
});

$('#url').tooltip({placement: 'bottom', title: 'Press CTRL+C to copy (or hold if mobile)', trigger: 'manual'});

function stats($http, $scope) {
    $http.get('/stats').success(function(data) {
        $scope.views = data.views.clean;
        $scope.links = data.links.clean;
    });
}

app.controller('pwCtrl', function($scope, $http) {
    $scope.$watch('decrypting', function() {
        if ($scope.decrypting) {
            $("#submit-btn").button("loading");
            $("#passwd").attr("disabled", "disabled");
        } else {
            $("#submit-btn").button("reset");
            $("#passwd").removeAttr("disabled");
        }
    });

    $scope.$watch('passwd', function() {
        // If the pw changes, disable the styles, we assume
        // they are fixing the pw
        $scope.error = false;
    });

    $scope.decrypt = function() {
        $scope.decrypting = true;
        $scope.error = false;

        $http.post('/decrypt', {path: $(location).attr('pathname'), password: $scope.passwd}).
          success(function(data) {
            // Redirect to the page
            window.location.replace(data.url);
          }).
          error(function(data) {
            $scope.decrypting = false;
            $scope.error = data.message ? data.message : "An unexpected error occured.";
            setTimeout(function(){$("#passwd").focus()}, 200);
          });
    };
});

app.controller('mainCtrl', function($scope, $http) {
    stats($http, $scope);
    $scope.$watch('shortening', function() {
        // Disable the shorten button and input box when we're
        // processing the url, otherwise re-enable them
        if ($scope.shortening) {
            $("#submit-btn").button("loading");
            $("#url").attr("disabled", "disabled");
        } else {
            $("#submit-btn").button("reset");
            $("#url").removeAttr("disabled");
        }
    });

    $scope.$watch('url', function() {
        // If the url changes, disable the styles, we assume
        // they are fixing the url, and hitting shorten again
        $scope.error = false;
        if (!$scope.finished) {
            $scope.success = false;
            $('#url').tooltip('hide');
            $("#submit-btn").removeAttr("disabled");
        }
        $scope.finished = false;
    });

    $scope.shorten = function() {
        var data, pw;
        $scope.shortening = true;
        $scope.error = false;
        $scope.success = false;

        // Add the password to the POST, if the user wants it
        data = $scope.passwd ? {url: $scope.url, password: $scope.passwd} : {url: $scope.url};

        $http.post('/add', data).
          success(function(data) {
            $scope.success = "Successfully shortened url!";
            $scope.shortening = false;
            // Use this, so the $watch.url doesn't override
            // our success styles
            $scope.finished = true;
            $scope.url = data.url;
            setTimeout(function(){
                $("#url").select();
                $('#url').tooltip('show');
                $("#submit-btn").attr("disabled", "disabled");
            },200);
          }).
          error(function(data) {
            $scope.shortening = false;
            $scope.success = false;
            $scope.error = data.message ? data.message : "An unexpected error occured.";
            setTimeout(function(){$("#url").focus()}, 200);
          });
    };
});
