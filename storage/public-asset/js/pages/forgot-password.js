function toastInfo(message) {
    Toastify({
        text: message,
        duration: 3000,
        close: true,
        gravity: "top",
        position: "right",
        style: {
            background: "#2c1c19"
        }
    }).showToast();
}

function ajaxLoading() {
    let form = document.getElementById("formForgot");

    // disable form elements
    for (let i = 0; i < form.elements.length; i++) {
        $(form.elements[i]).attr("disabled", true);
        if ($(form.elements[i]).attr("type") == "submit") {
            $(form.elements[i]).find("i").
                removeClass("fa-solid fa-envelope-open-text").
                addClass("spinner-border spinner-border-sm");
        }
    }

    $('#loginLink').on('click', function () {
        return false;
    }).addClass("text-muted");
}

function ajaxDone() {
    let form = document.getElementById("formForgot");

    // disable form login
    for (let i = 0; i < form.elements.length; i++) {
        $(form.elements[i]).removeAttr("disabled");
        if ($(form.elements[i]).attr("type") == "submit") {
            $(form.elements[i]).find("i").
                removeClass("spinner-border spinner-border-sm").
                addClass("fa-solid fa-envelope-open-text");
        }
    }

    $('#loginLink').off('click').removeClass("text-muted");
}

function submitForgot(form) {
    $.ajax({
        type: $(form).attr('method'),
        url: $(form).attr('action'),
        cache: false,
        beforeSend: function (xhr) {
            xhr.setRequestHeader('Accept', '*/*');
            ajaxLoading();
        },
        data: $(form).serialize(),
        dataType: "json"
    }).done(function (response) {
        toastInfo("email untuk reset password sudah terkirim")
        // go to submit token page
        setTimeout(function () {
            window.location.href = $('#loginLink').attr('href');
        }, 1000);
    }).fail(function ($jqXHR, $errorThrown) {
        ajaxDone();
        try {
            let response = JSON.parse($jqXHR.responseText);
            if (response.code) {
                switch (response.code) {
                    case 404002:
                        toastInfo("email tidak ditemukan!");
                        break;
                    case 400999:
                        // validator fails
                        for (let x in response.message) {
                            if (x == "email") {
                                for (let y in response.message[x]) {
                                    if (response.message[x][y].tag == "required") {
                                        toastInfo("email wajib diisi!");
                                    } else if (response.message[x][y].tag == "email") {
                                        toastInfo("email harus diisi dengan alamat email yang valid!");
                                    }
                                }
                            } else {
                                toastInfo("data yang diinput tidak valid");
                            }
                        }
                        break
                    default:
                        toastInfo("ada kesalahan teknis, error #" + response.code);
                        break;
                }
            }
        } catch (error) {
            toastInfo("ada kesalahan teknis");
            console.log(error);
        }
        return;
    });
}