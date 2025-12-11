// API Configuration
const API_BASE_URL = 'http://localhost:8080';

// State
let currentUser = null;
let currentProjectId = null;
let currentTaskId = null;
let chatPollInterval = null;
let isAdmin = false;

// ==================== API Functions ====================

// Auth API
const authAPI = {
    async register(fullName, email, password) {
        const response = await fetch(`${API_BASE_URL}/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ full_name: fullName, email, password }),
        });
        return response;
    },

    async login(email, password) {
        const response = await fetch(`${API_BASE_URL}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password }),
        });
        return response;
    },

    async verifyEmail(token) {
        const response = await fetch(`${API_BASE_URL}/auth/verify?token=${encodeURIComponent(token)}`, {
            method: 'GET',
        });
        return response;
    },
};

// Projects API
const projectsAPI = {
    async getProjects() {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_BASE_URL}/projects/`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
        return response;
    },

    async getProjectById(id) {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_BASE_URL}/projects/${id}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
        return response;
    },

    async createProject(title, description, status) {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_BASE_URL}/admin/projects/`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                title: title,
                description: description || '',
                status: status || 'active',
            }),
        });
        return response;
    },

    async updateProject(id, title, description, status) {
        const token = localStorage.getItem('token');
        const body = {};
        if (title !== null) body.title = title;
        if (description !== null) body.description = description;
        if (status !== null) body.status = status;

        const response = await fetch(`${API_BASE_URL}/admin/projects/${id}`, {
            method: 'PATCH',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(body),
        });
        return response;
    },

    async deleteProject(id) {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_BASE_URL}/admin/projects/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
        return response;
    },
};

// Tasks API
const tasksAPI = {
    async getTasks(projectId = null) {
        const token = localStorage.getItem('token');
        let url = `${API_BASE_URL}/tasks/`;
        if (projectId) {
            url += `?project_id=${projectId}`;
        }
        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
        return response;
    },

    async getTaskById(id) {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_BASE_URL}/tasks/${id}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
        return response;
    },

    async createTask(title, description, status, projectId, priority = 0) {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_BASE_URL}/admin/tasks/`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                title: title,
                description: description || '',
                status: status || 'todo',
                project_id: projectId,
                priority: priority,
            }),
        });
        return response;
    },

    async updateTask(id, title, description, status, priority) {
        const token = localStorage.getItem('token');
        const body = {};
        if (title !== null) body.title = title;
        if (description !== null) body.description = description;
        if (status !== null) body.status = status;
        if (priority !== null) body.priority = priority;

        const response = await fetch(`${API_BASE_URL}/admin/tasks/${id}`, {
            method: 'PATCH',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(body),
        });
        return response;
    },

    async deleteTask(id) {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_BASE_URL}/admin/tasks/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
        return response;
    },
};

// Chat API
const chatAPI = {
    async getMessages(type, id) {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_BASE_URL}/chat/${type}/${id}/`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
        return response;
    },

    async sendMessage(type, id, userId, text) {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_BASE_URL}/chat/${type}/${id}/`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                user_id: userId,
                text: text,
            }),
        });
        return response;
    },
};

// Users API
const usersAPI = {
    async getUserById(id) {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_BASE_URL}/users/${id}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
        return response;
    },

    async updateUser(id, fullName) {
        const token = localStorage.getItem('token');
        const body = {};
        if (fullName !== null) body.full_name = fullName;

        const response = await fetch(`${API_BASE_URL}/users/${id}`, {
            method: 'PATCH',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(body),
        });
        return response;
    },

    async deleteUser(id) {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_BASE_URL}/admin/users/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
        return response;
    },
};

// Reports API
const reportsAPI = {
    async getTopWorkers(projectId) {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_BASE_URL}/projects/${projectId}/reports/top-workers`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
        return response;
    },

    async getAverageTime(projectId) {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_BASE_URL}/projects/${projectId}/reports/avg-time`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
        return response;
    },

    async getCompletionPercent(projectId) {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_BASE_URL}/projects/${projectId}/reports/completion-percent`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
        return response;
    },

    async getUserTracker(projectId, userId) {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_BASE_URL}/projects/${projectId}/reports/user-tracker/${userId}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
        return response;
    },
};

// ==================== UI Functions ====================

function showError(elementId, message) {
    const element = document.getElementById(elementId);
    element.textContent = message;
    element.classList.add('show');
    setTimeout(() => {
        element.classList.remove('show');
    }, 5000);
}

function showSuccess(elementId, message) {
    const element = document.getElementById(elementId);
    element.textContent = message;
    element.classList.add('show');
    setTimeout(() => {
        element.classList.remove('show');
    }, 5000);
}

function clearErrors() {
    document.querySelectorAll('.error-message').forEach(el => {
        el.classList.remove('show');
        el.textContent = '';
    });
    document.querySelectorAll('.success-message').forEach(el => {
        el.classList.remove('show');
        el.textContent = '';
    });
}

// ==================== Auth Functions ====================

async function handleRegister(e) {
    e.preventDefault();
    clearErrors();

    const fullName = document.getElementById('register-fullname').value;
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;

    try {
        const response = await authAPI.register(fullName, email, password);
        const data = await response.json();

        if (response.ok) {
            showSuccess('register-success', 'Регистрация успешна! Проверьте вашу почту для подтверждения.');
            document.getElementById('register-form-element').reset();
            setTimeout(() => {
                document.getElementById('register-form').style.display = 'none';
                document.getElementById('verify-form').style.display = 'block';
            }, 2000);
        } else {
            showError('register-error', data.error || 'Ошибка при регистрации');
        }
    } catch (error) {
        console.error('Register error:', error);
        if (error.message && error.message.includes('Failed to fetch')) {
            showError('register-error', 'Не удалось подключиться к серверу. Убедитесь, что сервер запущен на http://localhost:8080');
        } else {
            showError('register-error', 'Ошибка соединения с сервером: ' + error.message);
        }
    }
}

async function handleLogin(e) {
    e.preventDefault();
    clearErrors();

    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;

    try {
        const response = await authAPI.login(email, password);
        const data = await response.json();

        if (response.ok) {
            localStorage.setItem('token', data.token);
            checkAdminStatus(); // Проверяем админ права после логина
            showApp();
            loadProjects();
        } else {
            showError('login-error', data.error || 'Ошибка при входе');
        }
    } catch (error) {
        console.error('Login error:', error);
        if (error.message && error.message.includes('Failed to fetch')) {
            showError('login-error', 'Не удалось подключиться к серверу. Убедитесь, что сервер запущен на http://localhost:8080');
        } else {
            showError('login-error', 'Ошибка соединения с сервером: ' + error.message);
        }
    }
}

async function handleVerifyEmail(e) {
    e.preventDefault();
    clearErrors();

    const token = document.getElementById('verify-token').value;

    if (!token) {
        showError('verify-error', 'Введите токен подтверждения');
        return;
    }

    try {
        const response = await authAPI.verifyEmail(token);
        const data = await response.json();

        if (response.ok) {
            showSuccess('verify-success', 'Email успешно подтвержден! Теперь вы можете войти.');
            document.getElementById('verify-form-element').reset();
            setTimeout(() => {
                document.getElementById('verify-form').style.display = 'none';
                document.getElementById('login-form').style.display = 'block';
            }, 2000);
        } else {
            showError('verify-error', data.error || 'Ошибка при подтверждении email');
        }
    } catch (error) {
        showError('verify-error', 'Ошибка соединения с сервером');
    }
}

function handleLogout() {
    localStorage.removeItem('token');
    currentUser = null;
    currentProjectId = null;
    currentTaskId = null;
    stopChatPolling();
    showAuth();
}

function showAuth() {
    document.getElementById('auth-screen').style.display = 'flex';
    document.getElementById('app-screen').style.display = 'none';
}

function showApp() {
    document.getElementById('auth-screen').style.display = 'none';
    document.getElementById('app-screen').style.display = 'block';
}

// ==================== Projects Functions ====================

async function loadProjects() {
    const projectsList = document.getElementById('projects-list');
    projectsList.innerHTML = '<div class="loading">Загрузка проектов...</div>';

    // Проверяем админ права и обновляем UI
    checkAdminStatus();
    updateAdminUI();

    try {
        const response = await projectsAPI.getProjects();
        const projects = await response.json();

        if (response.ok) {
            console.log('Projects loaded:', projects);
            if (projects.length === 0) {
                projectsList.innerHTML = '<div class="empty-state"><p>Проекты не найдены</p></div>';
            } else {
                projectsList.innerHTML = '';
                projects.forEach((project, index) => {
                    console.log(`Project ${index}:`, project, 'ID:', project.id || project.ID);
                    const projectCard = createProjectCard(project);
                    projectsList.appendChild(projectCard);
                });
            }
        } else {
            console.error('Failed to load projects:', response.status);
            projectsList.innerHTML = '<div class="empty-state"><p>Ошибка при загрузке проектов</p></div>';
        }
    } catch (error) {
        projectsList.innerHTML = '<div class="empty-state"><p>Ошибка соединения с сервером</p></div>';
    }
}

function createProjectCard(project) {
    const card = document.createElement('div');
    card.className = 'project-card';

    // Поддерживаем оба варианта имен полей (lowercase и с заглавной)
    // Проверяем ID более тщательно (может быть 0, что тоже валидно)
    let projectId = null;
    if (project.id !== undefined && project.id !== null) {
        projectId = project.id;
    } else if (project.ID !== undefined && project.ID !== null) {
        projectId = project.ID;
    }

    const projectTitle = project.title || project.Title || 'Без названия';
    const projectDesc = project.description || project.Description || 'Нет описания';
    const projectStatus = project.status || project.Status || 'N/A';

    console.log('Creating project card:', { project, projectId, title: projectTitle });

    if (projectId === null || projectId === undefined) {
        console.error('Project ID is missing in project object:', project);
    }

    card.innerHTML = `
        <h3>${escapeHtml(projectTitle)}</h3>
        <p>${escapeHtml(projectDesc)}</p>
        <span class="status-badge ${getStatusClass(projectStatus)}">${escapeHtml(projectStatus)}</span>
    `;
    card.addEventListener('click', () => {
        console.log('Opening project with ID:', projectId, 'from project:', project);
        if (projectId !== null && projectId !== undefined) {
            showProject(projectId);
        } else {
            console.error('Cannot open project: ID is missing!', project);
            alert('Ошибка: ID проекта не найден');
        }
    });
    return card;
}

async function showProject(projectId) {
    console.log('showProject called with ID:', projectId, 'type:', typeof projectId);

    // Более строгая проверка: ID должен быть числом и не null/undefined
    if (projectId === null || projectId === undefined || projectId === '') {
        console.error('Project ID is missing!', projectId);
        alert('Ошибка: ID проекта не указан');
        return;
    }

    // Преобразуем в число, если это строка
    const numericId = typeof projectId === 'string' ? parseInt(projectId, 10) : projectId;
    if (isNaN(numericId)) {
        console.error('Project ID is not a valid number!', projectId);
        alert('Ошибка: ID проекта должен быть числом');
        return;
    }

    projectId = numericId; // Используем числовой ID

    currentProjectId = projectId;
    currentTaskId = null;
    stopChatPolling();

    // Hide projects view, show project view
    document.getElementById('projects-view').style.display = 'none';
    document.getElementById('project-view').style.display = 'block';
    document.getElementById('task-view').style.display = 'none';

    // Проверяем админ права
    checkAdminStatus();
    const projectActions = document.getElementById('project-actions');
    const createTaskBtn = document.getElementById('create-task-btn');
    if (projectActions) projectActions.style.display = isAdmin ? 'flex' : 'none';
    if (createTaskBtn) createTaskBtn.style.display = isAdmin ? 'block' : 'none';

    // Показываем загрузку
    const titleEl = document.getElementById('project-title');
    const descEl = document.getElementById('project-description');
    const statusBadge = document.getElementById('project-status');
    const tasksList = document.getElementById('project-tasks-list');

    if (titleEl) titleEl.textContent = 'Загрузка...';
    if (descEl) descEl.textContent = 'Загрузка данных проекта...';
    if (statusBadge) statusBadge.textContent = '...';
    if (tasksList) tasksList.innerHTML = '<div class="loading">Загрузка задач...</div>';

    try {
        // Load project details
        console.log('Fetching project with ID:', projectId);
        const projectResponse = await projectsAPI.getProjectById(projectId);
        console.log('Project response status:', projectResponse.status, projectResponse.ok);

        if (!projectResponse.ok) {
            console.error('Failed to load project:', projectResponse.status, projectResponse.statusText);
            const errorText = await projectResponse.text();
            console.error('Error response:', errorText);
            let errorData = {};
            try {
                errorData = JSON.parse(errorText);
            } catch (e) {
                console.error('Failed to parse error response');
            }

            if (titleEl) titleEl.textContent = 'Ошибка загрузки';
            if (descEl) descEl.textContent = 'Не удалось загрузить данные проекта';
            if (statusBadge) statusBadge.textContent = 'N/A';
            if (tasksList) tasksList.innerHTML = '<div class="empty-state"><p>Ошибка при загрузке</p></div>';

            alert('Ошибка при загрузке проекта: ' + (errorData.error || projectResponse.statusText));
            return;
        }

        const project = await projectResponse.json();
        console.log('Project data loaded:', project);
        console.log('Project fields:', {
            title: project.title || project.Title,
            description: project.description || project.Description,
            status: project.status || project.Status
        });

        // Обновляем UI с данными проекта (поддерживаем оба варианта имен)
        if (titleEl) {
            const title = project.title || project.Title || 'Без названия';
            titleEl.textContent = title;
            console.log('Set title to:', title);
        }

        if (descEl) {
            const description = project.description || project.Description || 'Нет описания';
            descEl.textContent = description;
            console.log('Set description to:', description);
        }

        if (statusBadge) {
            const status = project.status || project.Status || 'N/A';
            statusBadge.textContent = status;
            statusBadge.className = `status-badge ${getStatusClass(status)}`;
            console.log('Set status to:', status);
        }

        // Load project tasks
        console.log('Fetching tasks for project ID:', projectId);
        const tasksResponse = await tasksAPI.getTasks(projectId);
        console.log('Tasks response status:', tasksResponse.status, tasksResponse.ok);

        const tasksList = document.getElementById('project-tasks-list');
        if (!tasksList) {
            console.error('Tasks list element not found!');
        }

        if (!tasksResponse.ok) {
            console.error('Failed to load tasks:', tasksResponse.status, tasksResponse.statusText);
            const errorText = await tasksResponse.text();
            console.error('Tasks error response:', errorText);

            if (tasksList) {
                tasksList.innerHTML = '<div class="empty-state"><p>Ошибка при загрузке задач</p></div>';
            }
        } else {
            const tasks = await tasksResponse.json();
            console.log('Tasks data loaded:', tasks);
            console.log('Tasks count:', Array.isArray(tasks) ? tasks.length : 'not an array');

            if (tasksList) {
                if (Array.isArray(tasks) && tasks.length > 0) {
                    tasksList.innerHTML = '';
                    tasks.forEach((task, index) => {
                        console.log(`Task ${index}:`, task);
                        const taskCard = createTaskCard(task);
                        tasksList.appendChild(taskCard);
                    });
                    console.log('Tasks rendered:', tasks.length);
                } else {
                    tasksList.innerHTML = '<div class="empty-state"><p>Задачи не найдены</p></div>';
                    console.log('No tasks found or tasks is not an array');
                }
            }
        }

        // Load reports
        console.log('Loading reports for project ID:', projectId);
        await loadReports(projectId);
    } catch (error) {
        console.error('Error loading project:', error);
        console.error('Error stack:', error.stack);

        const titleEl = document.getElementById('project-title');
        const descEl = document.getElementById('project-description');
        const statusBadge = document.getElementById('project-status');
        const tasksList = document.getElementById('project-tasks-list');

        if (titleEl) titleEl.textContent = 'Ошибка';
        if (descEl) descEl.textContent = 'Не удалось загрузить данные проекта: ' + error.message;
        if (statusBadge) statusBadge.textContent = 'N/A';
        if (tasksList) tasksList.innerHTML = '<div class="empty-state"><p>Ошибка при загрузке</p></div>';

        alert('Ошибка при загрузке проекта: ' + error.message);
    }
}

function createTaskCard(task) {
    const card = document.createElement('div');
    card.className = 'task-card';
    
    // Поддерживаем оба варианта имен полей
    const taskId = task.id || task.ID;
    const taskTitle = task.title || task.Title || 'Без названия';
    const taskDesc = task.description || task.Description || 'Нет описания';
    const taskStatus = task.status || task.Status || 'N/A';
    const taskPriority = task.priority || task.Priority || '';
    
    card.innerHTML = `
        <h4>${escapeHtml(taskTitle)}</h4>
        <p>${escapeHtml(taskDesc)}</p>
        <div class="task-meta">
            <span class="status-badge ${getStatusClass(taskStatus)}">${escapeHtml(taskStatus)}</span>
            ${taskPriority ? `<span>Приоритет: ${escapeHtml(taskPriority)}</span>` : ''}
        </div>
    `;
    card.addEventListener('click', () => {
        if (taskId) {
            showTask(taskId);
        } else {
            console.error('Task ID is missing!', task);
        }
    });
    return card;
}

// ==================== Task Functions ====================

async function showTask(taskId) {
    currentTaskId = taskId;
    stopChatPolling();

    // Hide project view, show task view
    document.getElementById('project-view').style.display = 'none';
    document.getElementById('task-view').style.display = 'block';

    // Проверяем админ права
    checkAdminStatus();
    const taskActions = document.getElementById('task-actions');
    if (taskActions) taskActions.style.display = isAdmin ? 'flex' : 'none';

    try {
        const response = await tasksAPI.getTaskById(taskId);
        const task = await response.json();

        if (response.ok) {
            document.getElementById('task-title').textContent = task.title || 'Без названия';
            document.getElementById('task-description').textContent = task.description || 'Нет описания';

            const statusBadge = document.getElementById('task-status');
            statusBadge.textContent = task.status || 'N/A';
            statusBadge.className = `status-badge ${getStatusClass(task.status)}`;

            document.getElementById('task-priority').textContent = task.priority !== undefined ? task.priority : 'N/A';

            // Load users
            const usersList = document.getElementById('task-users');
            if (task.users && task.users.length > 0) {
                usersList.innerHTML = '';
                task.users.forEach(user => {
                    const userBadge = document.createElement('span');
                    userBadge.className = 'user-badge';
                    userBadge.textContent = user.full_name || user.email || 'Unknown';
                    usersList.appendChild(userBadge);
                });
            } else {
                usersList.innerHTML = '<span class="empty-state">Нет участников</span>';
            }

            const chatInputContainer = document.getElementById('chat-input-container');
            const chatNoAccess = document.getElementById('chat-no-access');

            chatInputContainer.style.display = 'block';
            chatNoAccess.style.display = 'none';

            // Load chat messages
            await loadChatMessages('tasks', taskId);
            startChatPolling('tasks', taskId);
        }
    } catch (error) {
        console.error('Error loading task:', error);
    }
}

// ==================== Chat Functions ====================

async function loadChatMessages(type, id) {
    const messagesContainer = document.getElementById('chat-messages');

    try {
        const response = await chatAPI.getMessages(type, id);
        const messages = await response.json();

        if (response.ok) {
            if (messages.length === 0) {
                messagesContainer.innerHTML = '<div class="empty-state"><p>Нет сообщений</p></div>';
            } else {
                messagesContainer.innerHTML = '';
                messages.forEach(message => {
                    const messageElement = createChatMessage(message);
                    messagesContainer.appendChild(messageElement);
                });
                messagesContainer.scrollTop = messagesContainer.scrollHeight;
            }
        } else {
            messagesContainer.innerHTML = '<div class="empty-state"><p>Ошибка при загрузке сообщений</p></div>';
        }
    } catch (error) {
        console.error('Error loading chat messages:', error);
        messagesContainer.innerHTML = '<div class="empty-state"><p>Ошибка соединения с сервером</p></div>';
    }
}

async function createChatMessage(message) {
    // We need user info - let's try to get it
    let userName = 'Unknown';
    try {
        const userResponse = await usersAPI.getUserById(message.user_id);
        if (userResponse.ok) {
            const user = await userResponse.json();
            userName = user.full_name || user.email || 'Unknown';
        }
    } catch (error) {
        console.error('Error loading user:', error);
    }

    const messageDiv = document.createElement('div');
    messageDiv.className = 'chat-message';

    const createdAt = message.created_at ? new Date(message.created_at).toLocaleString('ru-RU') : 'N/A';

    messageDiv.innerHTML = `
        <div class="chat-message-header">
            <span class="chat-message-author">${escapeHtml(userName)}</span>
            <span class="chat-message-time">${createdAt}</span>
        </div>
        <div class="chat-message-text">${escapeHtml(message.text)}</div>
    `;

    return messageDiv;
}

async function handleSendMessage(e) {
    e.preventDefault();

    const text = document.getElementById('chat-input').value.trim();
    if (!text) return;

    // We need current user ID - since we don't have a /me endpoint,
    // we'll need to extract it from JWT token or store it after login
    // For now, let's try to get user ID from the task's users or store it
    // Actually, we need to decode JWT to get user ID
    // Let's create a simple JWT decoder (just for getting user_id)
    const token = localStorage.getItem('token');
    if (!token) {
        alert('Токен не найден. Пожалуйста, войдите снова.');
        return;
    }

    // Simple JWT decode (just for getting user_id from payload)
    // In production, use a proper JWT library
    let userId = null;
    try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        userId = payload.user_id || payload.userID;
    } catch (error) {
        console.error('Error decoding token:', error);
        alert('Ошибка при получении информации о пользователе');
        return;
    }

    if (!userId) {
        alert('Не удалось определить ID пользователя');
        return;
    }

    try {
        const response = await chatAPI.sendMessage('tasks', currentTaskId, userId, text);
        const data = await response.json();

        if (response.ok) {
            document.getElementById('chat-input').value = '';
            await loadChatMessages('tasks', currentTaskId);
        } else {
            if (response.status === 403) {
                // User doesn't have access
                document.getElementById('chat-input-container').style.display = 'none';
                document.getElementById('chat-no-access').style.display = 'block';
            } else {
                alert(data.error || 'Ошибка при отправке сообщения');
            }
        }
    } catch (error) {
        console.error('Error sending message:', error);
        alert('Ошибка соединения с сервером');
    }
}

function startChatPolling(type, id) {
    stopChatPolling();
    chatPollInterval = setInterval(() => {
        loadChatMessages(type, id);
    }, 3000); // Poll every 3 seconds
}

function stopChatPolling() {
    if (chatPollInterval) {
        clearInterval(chatPollInterval);
        chatPollInterval = null;
    }
}

// ==================== Utility Functions ====================

// ==================== Admin Functions ====================

function checkAdminStatus() {
    const token = localStorage.getItem('token');
    if (!token) {
        isAdmin = false;
        return false;
    }

    try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        isAdmin = payload.is_admin === true;
        return isAdmin;
    } catch (error) {
        console.error('Error decoding token:', error);
        isAdmin = false;
        return false;
    }
}

function updateAdminUI() {
    const createBtn = document.getElementById('create-project-btn');
    if (createBtn) {
        createBtn.style.display = isAdmin ? 'block' : 'none';
    }
}

// ==================== Project Creation Functions ====================

async function handleCreateProject(e) {
    e.preventDefault();
    clearErrors();

    const title = document.getElementById('project-title-input').value.trim();
    const description = document.getElementById('project-description-input').value.trim();
    const status = document.getElementById('project-status-input').value;

    if (!title) {
        showError('create-project-error', 'Название проекта обязательно');
        return;
    }

    try {
        const response = await projectsAPI.createProject(title, description, status);
        const data = await response.json();

        if (response.ok) {
            // Закрываем модальное окно
            document.getElementById('create-project-modal').style.display = 'none';
            document.getElementById('create-project-form').reset();

            // Обновляем список проектов
            loadProjects();
        } else {
            showError('create-project-error', data.error || 'Ошибка при создании проекта');
        }
    } catch (error) {
        console.error('Error creating project:', error);
        showError('create-project-error', 'Ошибка соединения с сервером');
    }
}

function openCreateProjectModal() {
    document.getElementById('create-project-modal').style.display = 'flex';
    document.getElementById('project-title-input').focus();
}

function closeCreateProjectModal() {
    document.getElementById('create-project-modal').style.display = 'none';
    document.getElementById('create-project-form').reset();
    clearErrors();
}

// ==================== Project Management Functions ====================

async function handleEditProject(e) {
    e.preventDefault();
    clearErrors();

    const title = document.getElementById('edit-project-title-input').value.trim();
    const description = document.getElementById('edit-project-description-input').value.trim();
    const status = document.getElementById('edit-project-status-input').value;

    if (!title) {
        showError('edit-project-error', 'Название проекта обязательно');
        return;
    }

    try {
        const response = await projectsAPI.updateProject(currentProjectId, title, description, status);
        const data = await response.json();

        if (response.ok) {
            closeEditProjectModal();
            showProject(currentProjectId); // Перезагружаем проект
        } else {
            showError('edit-project-error', data.error || 'Ошибка при обновлении проекта');
        }
    } catch (error) {
        console.error('Error updating project:', error);
        showError('edit-project-error', 'Ошибка соединения с сервером');
    }
}

async function handleDeleteProject() {
    if (!confirm('Вы уверены, что хотите удалить этот проект? Это действие нельзя отменить.')) {
        return;
    }

    try {
        const response = await projectsAPI.deleteProject(currentProjectId);
        const data = await response.json();

        if (response.ok) {
            backToProjects();
        } else {
            alert(data.error || 'Ошибка при удалении проекта');
        }
    } catch (error) {
        console.error('Error deleting project:', error);
        alert('Ошибка соединения с сервером');
    }
}

function openEditProjectModal() {
    // Заполняем форму текущими данными проекта
    const projectTitle = document.getElementById('project-title').textContent;
    const projectDescription = document.getElementById('project-description').textContent;
    const projectStatus = document.getElementById('project-status').textContent.toLowerCase();

    document.getElementById('edit-project-title-input').value = projectTitle;
    document.getElementById('edit-project-description-input').value = projectDescription;
    document.getElementById('edit-project-status-input').value = projectStatus;

    document.getElementById('edit-project-modal').style.display = 'flex';
}

function closeEditProjectModal() {
    document.getElementById('edit-project-modal').style.display = 'none';
    document.getElementById('edit-project-form').reset();
    clearErrors();
}

// ==================== Task Management Functions ====================

async function handleCreateTask(e) {
    e.preventDefault();
    clearErrors();

    const title = document.getElementById('task-title-input').value.trim();
    const description = document.getElementById('task-description-input').value.trim();
    const status = document.getElementById('task-status-input').value;
    const priority = parseInt(document.getElementById('task-priority-input').value) || 0;

    if (!title) {
        showError('create-task-error', 'Название задачи обязательно');
        return;
    }

    if (!currentProjectId) {
        showError('create-task-error', 'Проект не выбран');
        return;
    }

    try {
        console.log('Creating task:', { title, description, status, projectId: currentProjectId, priority });
        const response = await tasksAPI.createTask(title, description, status, currentProjectId, priority);
        console.log('Create task response status:', response.status, response.ok);

        if (!response.ok) {
            const errorText = await response.text();
            console.error('Failed to create task:', response.status, errorText);
            let errorData = {};
            try {
                errorData = JSON.parse(errorText);
            } catch (e) {
                console.error('Failed to parse error response');
            }
            showError('create-task-error', errorData.error || 'Ошибка при создании задачи');
            return;
        }

        const data = await response.json().catch(() => ({}));
        console.log('Task created successfully:', data);
        closeCreateTaskModal();
        showProject(currentProjectId); // Перезагружаем проект
    } catch (error) {
        console.error('Error creating task:', error);
        showError('create-task-error', 'Ошибка соединения с сервером: ' + error.message);
    }
}

async function handleEditTask(e) {
    e.preventDefault();
    clearErrors();

    const title = document.getElementById('edit-task-title-input').value.trim();
    const description = document.getElementById('edit-task-description-input').value.trim();
    const status = document.getElementById('edit-task-status-input').value;
    const priority = parseInt(document.getElementById('edit-task-priority-input').value) || 0;

    if (!title) {
        showError('edit-task-error', 'Название задачи обязательно');
        return;
    }

    try {
        const response = await tasksAPI.updateTask(currentTaskId, title, description, status, priority);
        const data = await response.json();

        if (response.ok) {
            closeEditTaskModal();
            showTask(currentTaskId); // Перезагружаем задачу
        } else {
            showError('edit-task-error', data.error || 'Ошибка при обновлении задачи');
        }
    } catch (error) {
        console.error('Error updating task:', error);
        showError('edit-task-error', 'Ошибка соединения с сервером');
    }
}

async function handleDeleteTask() {
    if (!confirm('Вы уверены, что хотите удалить эту задачу? Это действие нельзя отменить.')) {
        return;
    }

    try {
        const response = await tasksAPI.deleteTask(currentTaskId);
        const data = await response.json();

        if (response.ok) {
            if (currentProjectId) {
                showProject(currentProjectId); // Возвращаемся к проекту
            } else {
                backToProjects();
            }
        } else {
            alert(data.error || 'Ошибка при удалении задачи');
        }
    } catch (error) {
        console.error('Error deleting task:', error);
        alert('Ошибка соединения с сервером');
    }
}

function openCreateTaskModal() {
    document.getElementById('create-task-modal').style.display = 'flex';
    document.getElementById('task-title-input').focus();
}

function closeCreateTaskModal() {
    document.getElementById('create-task-modal').style.display = 'none';
    document.getElementById('create-task-form').reset();
    clearErrors();
}

function openEditTaskModal() {
    // Заполняем форму текущими данными задачи
    const taskTitle = document.getElementById('task-title').textContent;
    const taskDescription = document.getElementById('task-description').textContent;
    const taskStatus = document.getElementById('task-status').textContent.toLowerCase();
    const taskPriority = document.getElementById('task-priority').textContent;

    document.getElementById('edit-task-title-input').value = taskTitle;
    document.getElementById('edit-task-description-input').value = taskDescription;
    document.getElementById('edit-task-status-input').value = taskStatus === 'in progress' ? 'in_progress' : taskStatus;
    document.getElementById('edit-task-priority-input').value = taskPriority === 'N/A' ? 0 : taskPriority;

    document.getElementById('edit-task-modal').style.display = 'flex';
}

function closeEditTaskModal() {
    document.getElementById('edit-task-modal').style.display = 'none';
    document.getElementById('edit-task-form').reset();
    clearErrors();
}

// ==================== Utility Functions ====================

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function getStatusClass(status) {
    if (!status) return 'todo';
    const statusLower = status.toLowerCase();
    if (statusLower === 'in_progress' || statusLower === 'in-progress') return 'in_progress';
    if (statusLower === 'done' || statusLower === 'completed') return 'done';
    return 'todo';
}

// ==================== Navigation Functions ====================

function backToProjects() {
    currentProjectId = null;
    currentTaskId = null;
    stopChatPolling();
    document.getElementById('project-view').style.display = 'none';
    document.getElementById('task-view').style.display = 'none';
    document.getElementById('projects-view').style.display = 'block';
    loadProjects();
}

function backToProject() {
    currentTaskId = null;
    stopChatPolling();
    if (currentProjectId) {
        showProject(currentProjectId);
    } else {
        backToProjects();
    }
}

// ==================== Reports Functions ====================

async function loadReports(projectId) {
    const reportsContent = document.getElementById('reports-content');
    if (!reportsContent) return;

    reportsContent.innerHTML = '<div class="loading">Загрузка отчетов...</div>';

    try {
        const [topWorkersRes, avgTimeRes, completionRes] = await Promise.all([
            reportsAPI.getTopWorkers(projectId),
            reportsAPI.getAverageTime(projectId),
            reportsAPI.getCompletionPercent(projectId),
        ]);

        let html = '<div class="reports-grid">';

        // Top Workers
        if (topWorkersRes.ok) {
            const topWorkers = await topWorkersRes.json();
            html += '<div class="report-card"><h4>Топ работников</h4>';
            if (Array.isArray(topWorkers) && topWorkers.length > 0) {
                html += '<ul>';
                topWorkers.forEach(worker => {
                    html += `<li>${escapeHtml(worker.full_name || worker.email || 'Unknown')}: ${worker.task_count || 0} задач</li>`;
                });
                html += '</ul>';
            } else {
                html += '<p>Нет данных</p>';
            }
            html += '</div>';
        }

        // Average Time
        if (avgTimeRes.ok) {
            const avgTime = await avgTimeRes.json();
            html += '<div class="report-card"><h4>Среднее время</h4>';
            if (avgTime && avgTime.average_time) {
                html += `<p>${avgTime.average_time} дней</p>`;
            } else {
                html += '<p>Нет данных</p>';
            }
            html += '</div>';
        }

        // Completion Percent
        if (completionRes.ok) {
            const completion = await completionRes.json();
            html += '<div class="report-card"><h4>Процент завершения</h4>';
            if (completion && completion.completion_percent !== undefined) {
                html += `<p>${completion.completion_percent}%</p>`;
            } else {
                html += '<p>Нет данных</p>';
            }
            html += '</div>';
        }

        html += '</div>';
        reportsContent.innerHTML = html;
    } catch (error) {
        console.error('Error loading reports:', error);
        reportsContent.innerHTML = '<div class="empty-state"><p>Ошибка при загрузке отчетов</p></div>';
    }
}

// ==================== Event Listeners ====================

document.addEventListener('DOMContentLoaded', () => {
    // Check if user is already logged in
    const token = localStorage.getItem('token');
    if (token) {
        showApp();
        loadProjects();
    } else {
        showAuth();
    }

    // Auth form handlers
    document.getElementById('register-form-element').addEventListener('submit', handleRegister);
    document.getElementById('login-form-element').addEventListener('submit', handleLogin);
    document.getElementById('verify-form-element').addEventListener('submit', handleVerifyEmail);

    // Auth navigation
    document.getElementById('show-register').addEventListener('click', (e) => {
        e.preventDefault();
        clearErrors();
        document.getElementById('login-form').style.display = 'none';
        document.getElementById('register-form').style.display = 'block';
        document.getElementById('verify-form').style.display = 'none';
    });

    document.getElementById('show-login').addEventListener('click', (e) => {
        e.preventDefault();
        clearErrors();
        document.getElementById('register-form').style.display = 'none';
        document.getElementById('login-form').style.display = 'block';
        document.getElementById('verify-form').style.display = 'none';
    });

    document.getElementById('back-to-login').addEventListener('click', (e) => {
        e.preventDefault();
        clearErrors();
        document.getElementById('verify-form').style.display = 'none';
        document.getElementById('login-form').style.display = 'block';
    });

    // Logout
    document.getElementById('logout-btn').addEventListener('click', handleLogout);

    // Navigation
    document.getElementById('back-to-projects').addEventListener('click', backToProjects);
    document.getElementById('back-to-project').addEventListener('click', backToProject);

    // Chat form
    document.getElementById('chat-form').addEventListener('submit', handleSendMessage);

    // Create project handlers
    document.getElementById('create-project-btn').addEventListener('click', openCreateProjectModal);
    document.getElementById('close-create-modal').addEventListener('click', closeCreateProjectModal);
    document.getElementById('cancel-create-project').addEventListener('click', closeCreateProjectModal);
    document.getElementById('create-project-form').addEventListener('submit', handleCreateProject);

    // Закрытие модального окна при клике вне его
    document.getElementById('create-project-modal').addEventListener('click', (e) => {
        if (e.target.id === 'create-project-modal') {
            closeCreateProjectModal();
        }
    });

    // Edit project handlers
    const editProjectBtn = document.getElementById('edit-project-btn');
    const deleteProjectBtn = document.getElementById('delete-project-btn');
    if (editProjectBtn) {
        editProjectBtn.addEventListener('click', openEditProjectModal);
        document.getElementById('edit-project-form').addEventListener('submit', handleEditProject);
        document.getElementById('close-edit-project-modal').addEventListener('click', closeEditProjectModal);
        document.getElementById('cancel-edit-project').addEventListener('click', closeEditProjectModal);
        document.getElementById('edit-project-modal').addEventListener('click', (e) => {
            if (e.target.id === 'edit-project-modal') {
                closeEditProjectModal();
            }
        });
    }
    if (deleteProjectBtn) {
        deleteProjectBtn.addEventListener('click', handleDeleteProject);
    }

    // Create task handlers
    const createTaskBtn = document.getElementById('create-task-btn');
    if (createTaskBtn) {
        createTaskBtn.addEventListener('click', openCreateTaskModal);
        document.getElementById('create-task-form').addEventListener('submit', handleCreateTask);
        document.getElementById('close-create-task-modal').addEventListener('click', closeCreateTaskModal);
        document.getElementById('cancel-create-task').addEventListener('click', closeCreateTaskModal);
        document.getElementById('create-task-modal').addEventListener('click', (e) => {
            if (e.target.id === 'create-task-modal') {
                closeCreateTaskModal();
            }
        });
    }

    // Edit task handlers
    const editTaskBtn = document.getElementById('edit-task-btn');
    const deleteTaskBtn = document.getElementById('delete-task-btn');
    if (editTaskBtn) {
        editTaskBtn.addEventListener('click', openEditTaskModal);
        document.getElementById('edit-task-form').addEventListener('submit', handleEditTask);
        document.getElementById('close-edit-task-modal').addEventListener('click', closeEditTaskModal);
        document.getElementById('cancel-edit-task').addEventListener('click', closeEditTaskModal);
        document.getElementById('edit-task-modal').addEventListener('click', (e) => {
            if (e.target.id === 'edit-task-modal') {
                closeEditTaskModal();
            }
        });
    }
    if (deleteTaskBtn) {
        deleteTaskBtn.addEventListener('click', handleDeleteTask);
    }
});

