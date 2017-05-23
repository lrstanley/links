// Initialize all callbacks.
$(function () {
    $("#url_box").focus();
    $('[data-toggle="tooltip"]').tooltip();
    var clipboard = new Clipboard('.clip');

    clipboard.on('success', function (e) {
        notie.alert({time: 1, text: "Copied to clipboard"})
    });
});
