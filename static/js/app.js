document.addEventListener('DOMContentLoaded', function() {
    const loginContainer = document.getElementById('loginContainer');
    const registerContainer = document.getElementById('registerContainer');
    const loginButton = document.getElementById('loginButton');
    const registerButton = document.getElementById('registerButton');
    const overlay = document.createElement('div');
    overlay.classList.add('overlay');

    function toggleOverlay() {
        overlay.classList.toggle('show-overlay');
    }

    function toggleLogin() {
        loginContainer.classList.toggle('show-login-register');
    }

    function toggleRegister() {
        registerContainer.classList.toggle('show-login-register');
    }

    function handleClickOutside(event) {
        if (!loginContainer.contains(event.target) && !loginButton.contains(event.target) &&
            !registerContainer.contains(event.target) && !registerButton.contains(event.target)) {
            loginContainer.classList.remove('show-login-register');
            registerContainer.classList.remove('show-login-register');
            overlay.classList.remove('show-overlay');
        }
    }

    if (loginButton) {
        loginButton.addEventListener('click', function(event) {
            event.preventDefault();
            toggleOverlay();
            toggleLogin();
        });
    }

    if (registerButton) {
        registerButton.addEventListener('click', function(event) {
            event.preventDefault();
            toggleOverlay();
            toggleRegister();
        });
    }

    document.body.appendChild(overlay);
    document.body.addEventListener('click', handleClickOutside);
});
