document.getElementById("filterForm").addEventListener("submit", function(event) {
    event.preventDefault(); // Prevent the default form submission

    const form = event.target;
    const formData = new FormData(form);
    const params = new URLSearchParams();

    formData.forEach((value, key) => {
        params.append(key, value);
    });

    // Redirect to the URL with query parameters
    window.location.href = `/?${params.toString()}`;
});
