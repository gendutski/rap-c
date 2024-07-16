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

function formLoginLoading() {
    let formLogin = document.getElementById("formLogin");
    let formGuest = document.getElementById("formGuest");

    // disable form login
    for (let i = 0; i < formLogin.elements.length; i++) {
        $(formLogin.elements[i]).attr("disabled", true);
        if ($(formLogin.elements[i]).attr("type") == "submit") {
            $(formLogin.elements[i]).find("i").
                removeClass("fa-solid fa-door-open").
                addClass("spinner-border spinner-border-sm");
        }
    }

    // disable form guest
    for (let i = 0; i < formGuest.elements.length; i++) {
        $(formGuest.elements[i]).attr("disabled", true);
    }

    // disable forgot password
    $('#forgotLink').on("click", function () {
        return false;
    }).addClass("text-muted");
}

function formLoginDone() {
    let formLogin = document.getElementById("formLogin");
    let formGuest = document.getElementById("formGuest");

    // disable form login
    for (let i = 0; i < formLogin.elements.length; i++) {
        $(formLogin.elements[i]).removeAttr("disabled");
        if ($(formLogin.elements[i]).attr("type") == "submit") {
            $(formLogin.elements[i]).find("i").
                removeClass("spinner-border spinner-border-sm").
                addClass("fa-solid fa-door-open");
        }
    }

    // disable form guest
    for (let i = 0; i < formGuest.elements.length; i++) {
        $(formGuest.elements[i]).removeAttr("disabled");
    }

    // enable forgot password
    $('#forgotLink').off("click").removeClass("text-muted");
}

function formGuestLoading() {
    let formLogin = document.getElementById("formLogin");
    let formGuest = document.getElementById("formGuest");

    // disable form login
    for (let i = 0; i < formLogin.elements.length; i++) {
        $(formLogin.elements[i]).attr("disabled", true);
    }

    // disable form guest
    for (let i = 0; i < formGuest.elements.length; i++) {
        $(formGuest.elements[i]).attr("disabled", true);
        if ($(formGuest.elements[i]).attr("type") == "submit") {
            $(formGuest.elements[i]).find("i").
                removeClass("fa-solid fa-masks-theater").
                addClass("spinner-border spinner-border-sm");
        }
    }

    // disable forgot password
    $('#forgotLink').on("click", function () {
        return false;
    }).addClass("text-muted");
}

function formGuestDone() {
    let formLogin = document.getElementById("formLogin");
    let formGuest = document.getElementById("formGuest");

    // disable form login
    for (let i = 0; i < formLogin.elements.length; i++) {
        $(formLogin.elements[i]).removeAttr("disabled");
    }

    // disable form guest
    for (let i = 0; i < formGuest.elements.length; i++) {
        $(formGuest.elements[i]).removeAttr("disabled");
        if ($(formLogin.elements[i]).attr("type") == "submit") {
            $(formLogin.elements[i]).find("i").
                removeClass("spinner-border spinner-border-sm").
                addClass("fa-solid fa-masks-theater");
        }
    }

    // enable forgot password
    $('#forgotLink').off("click").removeClass("text-muted");
}

function submitLogin(form) {
    $.ajax({
        type: $(form).attr('method'),
        url: $(form).attr('action'),
        cache: false,
        beforeSend: function (xhr) {
            xhr.setRequestHeader('Accept', '*/*');
            formLoginLoading();
        },
        data: $(form).serialize(),
        dataType: "json"
    }).done(function (response) {
        // go to submit token page
        $('#formSubmitToken input[name="token"]').val(response.token);
        $('#formSubmitToken').submit();
    }).fail(function ($jqXHR) {
        formLoginDone();
        try {
            let response = JSON.parse($jqXHR.responseText);
            if (response.code) {
                switch (response.code) {
                    case 401001:
                        toastInfo("email atau password salah");
                        break;
                    case 401002:
                        toastInfo("user non aktif tidak dapat login");
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

function submitGuest(form) {
    $.ajax({
        type: $(form).attr('method'),
        url: $(form).attr('action'),
        cache: false,
        beforeSend: function (xhr) {
            xhr.setRequestHeader('Accept', '*/*');
            formGuestLoading();
        },
        data: $(form).serialize(),
        dataType: "json"
    }).done(function (response) {
        // go to submit token page
        $('#formSubmitToken input[name="token"]').val(response.token);
        $('#formSubmitToken').submit();
    }).fail(function ($jqXHR, $errorThrown) {
        formGuestDone();
        try {
            let response = JSON.parse($jqXHR.responseText);
            if (response.code) {
                switch (response.code) {
                    case 401001:
                        toastInfo("email atau password salah");
                        break;
                    case 401002:
                        toastInfo("user non aktif tidak dapat login");
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