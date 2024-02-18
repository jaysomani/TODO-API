// Function to fetch all TODO items from the server and display them on the webpage
async function displayAllTodos() {
    try {
        const response = await fetch('/todos/read');
        const todos = await response.json();
        
        const todoList = document.getElementById('todo-list');
        todoList.innerHTML = ''; // Clear the existing todo items
        
        todos.forEach(todo => {
            const todoItem = document.createElement('li');
            todoItem.classList.add('todo-item');
            
            const todoContent = document.createElement('div');
            todoContent.textContent = `${todo.title} - ${todo.description}`;

            const statusIndicator = document.createElement('span');
            statusIndicator.textContent = todo.status;
            statusIndicator.classList.add('status-indicator');
            statusIndicator.style.backgroundColor = todo.status === 'completed' ? '#28a745' : '#dc3545';

            const editButton = document.createElement('button');
            editButton.textContent = '✏️';
            editButton.classList.add('edit-btn');
            editButton.addEventListener('click', () => {
                editTodoItem(todo);
            });
            
            const deleteButton = document.createElement('button');
            deleteButton.textContent = '❌';
            deleteButton.classList.add('delete-btn');
            deleteButton.addEventListener('click', async () => {
                try {
                    const response = await fetch('/todos/delete', {
                        method: 'DELETE',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify(todo.id) 
                    });
                    if (response.ok) {
                        todoList.removeChild(todoItem);
                    } else {
                        throw new Error('Failed to delete todo item');
                    }
                } catch (error) {
                    console.error('Error deleting todo item:', error);
                    alert('Failed to delete todo item. Please try again.');
                }
            });
            const statusButton = document.createElement('button');
            statusButton.textContent = todo.status === 'completed' ? 'Pending' : 'Complete';
            statusButton.classList.add('status-btn');
            statusButton.addEventListener('click', async () => {
                todo.status = todo.status === 'completed' ? 'pending' : 'completed';
                const response = await fetch('/todos/update', {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(todo)
                });
                if (response.ok) {
                    statusIndicator.textContent = todo.status;
                    statusIndicator.style.backgroundColor = todo.status === 'completed' ? '#28a745' : '#dc3545';
                    statusButton.textContent = todo.status === 'completed' ? 'Pending' : 'Complete';
                } else {
                    console.error('Error updating todo status');
                    alert('Failed to update todo status. Please try again.');
                }
            });
            todoItem.appendChild(todoContent);
            todoItem.appendChild(statusIndicator);
            todoItem.appendChild(editButton);
            todoItem.appendChild(deleteButton);
            todoItem.appendChild(statusButton);
            todoList.appendChild(todoItem);
        });
    } catch (error) {
        console.error('Error fetching todos:', error);
    }
}

// Function to handle editing a todo item
function editTodoItem(todo) {
    const newTitle = prompt('Enter new title:', todo.title);
    const newDescription = prompt('Enter new description:', todo.description);
    
    if (newTitle !== null && newDescription !== null) {
        const updatedTodo = {
            ...todo,
            title: newTitle,
            description: newDescription
        };
        
        fetch('/todos/update', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(updatedTodo)
        })
        .then(response => {
            if (response.ok) {
                displayAllTodos();
                alert('Todo item updated successfully!');
            } else {
                throw new Error('Failed to update todo item');
            }
        })
        .catch(error => {
            console.error('Error updating todo item:', error);
            alert('Failed to update todo item. Please try again.');
        });
    }
}

// Call the displayAllTodos function when the page loads
window.onload = displayAllTodos;

// Event listener for Add Todo button
document.getElementById('btn-create').addEventListener('click', async () => {
    const title = document.getElementById('todo-title').value;
    const description = document.getElementById('todo-description').value;
    try {
        const response = await fetch('/todos', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ title, description, status: 'pending' }) // Include status in the request body
        });
        if (response.ok) {
            alert('Todo item created successfully!');
            displayAllTodos();
        } else {
            throw new Error('Failed to add todo item');
        }
    } catch (error) {
        console.error('Error adding todo item:', error);
        alert('Failed to add todo item. Please try again.');
    }
});