// Initialize all callbacks.
$(function () {
    $(".input-focus").focus();
    var clipboard = new Clipboard('.clip');

    clipboard.on('success', function (e) {
        notie.alert({time: 1, text: "Copied to clipboard"})
    });
});
