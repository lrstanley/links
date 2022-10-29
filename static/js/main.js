/**
 * Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
 * of this source code is governed by the MIT license that can be found in
 * the LICENSE file.
 */

// Initialize all callbacks.
$(function () {
    $(".input-focus").focus();
    var clipboard = new Clipboard('.clip');

    clipboard.on('success', function (e) {
        notie.alert({time: 1, text: "Copied to clipboard"})
    });
});
