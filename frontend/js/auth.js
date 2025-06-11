// Show login form as modal
function showLoginForm() {
    document.getElementById('login-form').style.display = 'flex';
    document.getElementById('signup-form').style.display = 'none';
}

// Show signup form as modal
function showSignupForm() {
    document.getElementById('signup-form').style.display = 'flex';
    document.getElementById('login-form').style.display = 'none';
}

// Close login form
function closeL() {
    document.getElementById('login-form').style.display = 'none';
}

// Close signup form
function closeS() {
    document.getElementById('signup-form').style.display = 'none';
}

// Close forms when clicking outside
document.getElementById('login-form').addEventListener('click', function(e) {
    if (e.target === this) {
        closeL();
    }
});

document.getElementById('signup-form').addEventListener('click', function(e) {
    if (e.target === this) {
        closeS();
    }
});

// Login form handler
document.getElementById('loginForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const identifier = document.querySelector('#loginForm input[name="identifier"]').value;
    const password = document.querySelector('#loginForm input[name="password"]').value;

    try {
        const response = await fetch('/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                identifier,
                password,
            }),
        });

        const data = await response.text();
        const errorDiv = document.getElementById('login-error');

        if (response.ok) {
            window.location.href = '/';
        } else {
            errorDiv.textContent = data;
            errorDiv.style.display = 'block';
        }
    } catch (error) {
        console.error('Error:', error);
        document.getElementById('login-error').textContent = 'An error occurred. Please try again.';
        document.getElementById('login-error').style.display = 'block';
    }
});

// Signup form handler
document.getElementById('signupForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const form = e.target;
    const data = {
        nickname: form.nickname.value,
        firstName: form.firstName.value,
        lastName: form.lastName.value,
        age: parseInt(form.age.value, 10),
        gender: form.gender.value,
        email: form.email.value,
        password: form.password.value,
    };

    try {
        const response = await fetch('/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });

        const errorDiv = document.getElementById('signup-error');
        if (response.ok) {
            errorDiv.textContent = 'Registration successful! Redirecting...';
            errorDiv.style.display = 'block';
            setTimeout(() => {
                showLoginForm();
            }, 2000);
        } else {
            const errorText = await response.text();
            errorDiv.textContent = 'Error: ' + errorText;
            errorDiv.style.display = 'block';
        }
    } catch (error) {
        console.error('Error:', error);
        document.getElementById('signup-error').textContent = 'An error occurred. Please try again.';
        document.getElementById('signup-error').style.display = 'block';
    }
});

// Show password functionality
document.querySelectorAll('.icon').forEach(icon => {
    icon.addEventListener('click', function() {
        const passwordField = this.previousElementSibling;
        if (passwordField.type === 'password') {
            passwordField.type = 'text';
            this.textContent = '🔒';
        } else {
            passwordField.type = 'password';
            this.textContent = '👁️';
        }
    });
});

// Logo click to return to homepage
document.querySelector('.logo').addEventListener('click', function() {
    showHomepage();
});
