
        // Handle post creation form submission
        document.getElementById('createPostForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const formData = new FormData(e.target);
            const response = await fetch('/post', {
                method: 'POST',
                body: formData
            });

            if (response.ok) {
                alert('Post created successfully!');
                window.location.href = '/';
            } else if (response.status === 401) {
                alert('Failed to create post: user not logged in.');
            } else {
                alert('Failed to create post.');
            }
        });

        // Handle like/dislike form submissions
        // Like/dislike handler for both posts and comments
        document.addEventListener('DOMContentLoaded', () => {
            document.querySelectorAll('.like-form').forEach(form => {
                form.addEventListener('submit', async (event) => {
                    event.preventDefault();

                    const formData = new URLSearchParams();
                    formData.append('id', event.target.querySelector('input[name="id"]').value);
                    formData.append('item_type', event.target.querySelector('input[name="item_type"]').value);
                    formData.append('type', event.submitter.value);

                    try {
                        const response = await fetch('/likes', {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/x-www-form-urlencoded',
                            },
                            body: formData
                        });

                        if (response.ok) {
                            //error handling for easy debugging of the code
                            const itemType = event.target.querySelector('input[name="item_type"]').value;
                            const itemId = event.target.querySelector('input[name="id"]').value;
                            const countSpan = document.getElementById(`${event.submitter.value}-count-${itemType}-${itemId}`);

                            if (countSpan) {
                                const currentCount = parseInt(countSpan.textContent);
                                // countSpan.textContent = currentCount + 1;
                                window.location.reload()
                            }
                        } else if (response.status === 401) {
                            window.location.href = '/login';
                        } else {
                            const errorText = await response.text();
                            console.error('Error:', errorText);
                            alert('Failed to process like/dislike');
                        }
                    } catch (error) {
                        console.error('Request failed:', error);
                        alert('Failed to process request');
                    }
                });
            });
        });

        // Handle comment form submissions
        document.querySelectorAll('.comment-form').forEach(form => {
            form.addEventListener('submit', async (event) => {
                event.preventDefault();

                const formData = new URLSearchParams(new FormData(event.target));

                try {
                    const response = await fetch('/add_comment', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/x-www-form-urlencoded',
                        },
                        body: formData
                    });

                    if (response.ok) {
                        alert('Comment created successfully');
                        window.location.reload();
                    } else if (response.status === 401) {
                        alert('Failed to comment: user not logged in.');
                        window.location.href = '/login';
                    } else {
                        alert('Failed to comment.');
                    }
                } catch (error) {
                    console.error('Request failed:', error);
                }
            });
        });