function decrypt() {
    // Add the password to the POST, is the user wants it
    if ($("#encryption").val().length == 0) {
        // Empty input box
        $("#password").removeClass("has-success");
        $("#password").addClass("has-error");
        // Add an alert with the success/failure information
        $("#messages").removeClass();
        $("#messages").addClass("alert alert-danger");
        $("#messages").html("<strong>Please enter a password<strong>");
        return
    }
    // Assume it has a password, but we don't know if it's wrong/right yet
    data = {path: $(location).attr('pathname'), password: $("#encryption").val()};
    // POST data to the server to check...
    $.post("decrypt", data).done(function(data) {
        var data = $.parseJSON(data);
        if (data["success"]) {
            // If success, redirect the user to the URL
            // similar behavior as an HTTP redirect
            window.location.replace(data["url"]);
            return
        } else {
            // Change form validation state here
            $("#password").removeClass("has-success");
            $("#password").addClass("has-error");

            // Add an alert with the success/failure information
            $("#messages").removeClass();
            $("#messages").addClass("alert alert-danger");
            $("#messages").html("<strong>" + data["message"] + "<strong>");
            $("#encryption").val('');
        }
        $("#submit-btn").button("reset");
    });
}

function add() {
    $("#submit-btn").button("loading");
    $("#url").attr("disabled", "disabled");
    // Add the password to the POST, if the user wants it
    if ($("#password-box").val().length == 0) {
        data = {url: $("#url").val()};
    } else {
        data = {url: $("#url").val(), password: $("#password-box").val()};
    }
    $.post("add", data).done(function(data) {
        var data = $.parseJSON(data);
        if (data["success"]) {
            // Change form validation state here
            $("#shortener").removeClass("has-error");
            $("#shortener").addClass("has-success");

            // Clear out the old input boxs
            $("#url").val('');
            $("#password-box").val('');

            // Add an alert with the success/failure information
            $("#messages").removeClass();
            $("#messages").addClass("alert alert-success");
            $("#messages").html("<strong>" + 'Successfully shortened your URL!' +
                '<input type="text" value="' + data["url"] + '" class="form-control tip"' +
                'id="url-to-copy" title="Press CTRL+C to copy" readonly>' + "<strong>");

            // Specifically highlight the URL input-box, and select it
            $('#url-to-copy').focus();
            $('#url-to-copy').select();
            $('#url-to-copy').tooltip("show");
        } else {
            // Change form validation state here
            $("#shortener").removeClass("has-success");
            $("#shortener").addClass("has-error");

            // Add an alert with the success/failure information
            $("#messages").removeClass();
            $("#messages").addClass("alert alert-danger");
            $("#messages").html("<strong>" + data["message"] + "<strong>");
        }
        $("#submit-btn").button("reset");
        $("#url").removeAttr("disabled");
    });
}

$("input").focus(function() {
  this.value = "";
});


$(function () {
   $("#url").focus();
   $("#encryption").focus();
});

$('#url').keypress(function(e) {
    if (e.which == 13) {
        add();
        return false;
    }
});

$(function(){
    $('.tip').tooltip({
        placement: "top",
        trigger: "manual"
    });
});