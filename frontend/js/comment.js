// Comment section toggle
document.addEventListener('click', (e) => {
    if (e.target.classList.contains('view-comments')) {
        const postId = e.target.dataset.postId;
        const commentsSection = document.getElementById(`comments-${postId}`);
        commentsSection.style.display = commentsSection.style.display === 'none' ? 'block' : 'none';
    }
});

// Comment form handling
document.addEventListener('submit', async (e) => {
    if (e.target.classList.contains('comment-form')) {
        e.preventDefault();

        const form = e.target;
        const submitButton = form.querySelector('button[type="submit"]');
        const originalButtonText = submitButton.textContent;
        submitButton.disabled = true;
        submitButton.textContent = 'Submitting...';

        const postId = form.dataset.postId;
        const content = form.querySelector('textarea').value;

        if (!postId || !content.trim()) {
            alert('Comment content cannot be empty.');
            submitButton.disabled = false;
            submitButton.textContent = originalButtonText;
            return;
        }
        
        try {
            const response = await fetch("/api/comment", {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ postId: Number(postId), content: content })
            });
             
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`Failed to create comment: ${errorText}`);
            }

            // Show success message
            const successMessage = document.createElement('div');
            successMessage.className = 'success-message';
            successMessage.textContent = 'Comment added successfully!';
            form.insertBefore(successMessage, form.firstChild);
            
            // Clear form and reload after a short delay
            setTimeout(() => {
                form.reset();
                window.location.reload();
            }, 1500);
        } catch (error) {
            console.error('Error creating comment:', error);
            alert('Failed to create comment. Please try again.');
        } finally {
            submitButton.disabled = false;
            submitButton.textContent = originalButtonText;
        }
    }
});
