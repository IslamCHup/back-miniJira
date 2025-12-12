// API Configuration
const API_BASE_URL = 'http://localhost:8080';

// State
let currentUser = null;
let currentProjectId = null;
let currentTaskId = null;
let chatPollInterval = null;
let isAdmin = false;
const userCache = {}; // Кэш для пользователей
let currentTaskData = null; // Храним исходные данные текущей задачи
let originalTaskStatus = null; // Храним исходный статус задачи для сравнения

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
    async getTasks(projectId = null, filters = {}) {
        const token = localStorage.getItem('token');
        let url = `${API_BASE_URL}/tasks/`;
        const params = [];

        if (projectId) {
            params.push(`project_id=${projectId}`);
        }

        // Добавляем фильтры как query параметры
        if (filters.status) {
            params.push(`status=${encodeURIComponent(filters.status)}`);
        }
        if (filters.user_id) {
            params.push(`user_id=${filters.user_id}`);
        }
        if (filters.search) {
            params.push(`search=${encodeURIComponent(filters.search)}`);
        }
        if (filters.priority !== undefined && filters.priority !== null && filters.priority !== '') {
            params.push(`priority=${filters.priority}`);
        }
        if (filters.sort_by) {
            params.push(`sort_by=${encodeURIComponent(filters.sort_by)}`);
        }
        if (filters.sort_order) {
            params.push(`sort_order=${encodeURIComponent(filters.sort_order)}`);
        }

        if (params.length > 0) {
            url += '?' + params.join('&');
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

        // Отправляем только непустые значения
        if (title !== null && title !== undefined && title.trim() !== '') {
            body.title = title.trim();
        }
        if (description !== null && description !== undefined) {
            body.description = description.trim();
        }
        if (status !== null && status !== undefined && status !== '') {
            body.status = status;
        }
        if (priority !== null && priority !== undefined) {
            body.priority = priority;
        }

        console.log('Update task request body:', body);

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

        // Load project tasks with filters
        console.log('Fetching tasks for project ID:', projectId);
        const taskFilters = getTaskFilters();
        const tasksResponse = await tasksAPI.getTasks(projectId, taskFilters);
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
    console.log('showTask called with ID:', taskId, 'type:', typeof taskId);

    if (!taskId) {
        console.error('Task ID is missing!');
        alert('Ошибка: ID задачи не указан');
        return;
    }

    // Преобразуем в число, если это строка
    const numericId = typeof taskId === 'string' ? parseInt(taskId, 10) : taskId;
    if (isNaN(numericId)) {
        console.error('Task ID is not a valid number!', taskId);
        alert('Ошибка: ID задачи должен быть числом');
        return;
    }

    currentTaskId = numericId;
    console.log('Current task ID set to:', currentTaskId);
    stopChatPolling();

    // Hide project view, show task view
    document.getElementById('project-view').style.display = 'none';
    document.getElementById('task-view').style.display = 'block';

    // Проверяем админ права
    checkAdminStatus();
    const taskActions = document.getElementById('task-actions');
    if (taskActions) taskActions.style.display = isAdmin ? 'flex' : 'none';

    try {
        console.log('Fetching task with ID:', numericId);
        const response = await tasksAPI.getTaskById(numericId);
        console.log('Task response status:', response.status, response.ok);

        if (!response.ok) {
            const errorText = await response.text().catch(() => 'Unknown error');
            console.error('Failed to load task:', response.status, errorText);
            alert('Ошибка при загрузке задачи');
            return;
        }

        const task = await response.json();
        console.log('Task data loaded:', task);

        if (response.ok) {
            // Сохраняем исходные данные задачи для редактирования
            currentTaskData = task;

            // Сохраняем исходный статус задачи для сравнения при редактировании
            if (task.status || task.Status) {
                let taskStatus = (task.status || task.Status).toLowerCase();
                if (taskStatus === 'in progress' || taskStatus === 'in-progress') {
                    taskStatus = 'in_progress';
                }
                originalTaskStatus = taskStatus;
            }

            document.getElementById('task-title').textContent = task.title || task.Title || 'Без названия';
            document.getElementById('task-description').textContent = task.description || task.Description || 'Нет описания';

            const statusBadge = document.getElementById('task-status');
            const taskStatus = task.status || task.Status || 'N/A';
            statusBadge.textContent = taskStatus;
            statusBadge.className = `status-badge ${getStatusClass(taskStatus)}`;

            // Отображаем приоритет
            const priorityEl = document.getElementById('task-priority');
            if (priorityEl) {
                const priorityText = task.priority || task.Priority || 'N/A';
                priorityEl.textContent = priorityText;

                // Определяем исходное числовое значение приоритета из текста
                let priorityNum = 0;
                if (priorityText && priorityText !== 'N/A') {
                    if (priorityText.includes('Очень важно')) {
                        priorityNum = 2;
                    } else if (priorityText.includes('Важно')) {
                        priorityNum = 1;
                    } else {
                        // Пытаемся извлечь число из строки
                        const numMatch = priorityText.match(/\d+/);
                        if (numMatch) {
                            priorityNum = parseInt(numMatch[0]);
                        }
                    }
                }
                // Сохраняем в data-атрибуте для использования при редактировании
                priorityEl.setAttribute('data-priority', priorityNum);
                console.log('Task priority saved:', priorityNum, 'from text:', priorityText);
            }

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

            if (chatInputContainer) chatInputContainer.style.display = 'block';
            if (chatNoAccess) chatNoAccess.style.display = 'none';

            // Load chat messages (используем numericId, который мы установили в currentTaskId)
            console.log('Loading chat for task ID:', currentTaskId);
            await loadChatMessages('tasks', currentTaskId);
            startChatPolling('tasks', currentTaskId);
        }
    } catch (error) {
        console.error('Error loading task:', error);
    }
}

// ==================== Chat Functions ====================

async function loadChatMessages(type, id) {
    const messagesContainer = document.getElementById('chat-messages');
    if (!messagesContainer) {
        console.error('Chat messages container not found');
        return;
    }

    if (!id) {
        console.error('Chat ID is missing');
        messagesContainer.innerHTML = '<div class="empty-state"><p>ID чата не указан</p></div>';
        return;
    }

    try {
        console.log('Loading chat messages:', { type, id });
        const response = await chatAPI.getMessages(type, id);
        console.log('Chat messages response status:', response.status, response.ok);
        console.log('Response headers:', Object.fromEntries(response.headers.entries()));

        if (!response.ok) {
            let errorText = 'Unknown error';
            try {
                errorText = await response.text();
            } catch (e) {
                console.error('Failed to read error response:', e);
            }
            console.error('Failed to load chat messages:', response.status, errorText);

            let errorMessage = 'Ошибка при загрузке сообщений';
            try {
                const errorData = JSON.parse(errorText);
                if (errorData.error) {
                    errorMessage = errorData.error;
                }
            } catch (e) {
                // Если не JSON, используем текст как есть
                if (errorText && errorText !== 'Unknown error') {
                    errorMessage = errorText;
                }
            }

            messagesContainer.innerHTML = `<div class="empty-state"><p>${escapeHtml(errorMessage)}</p></div>`;
            return;
        }

        // Парсим JSON ответ
        let messages = [];
        try {
            messages = await response.json();
            console.log('Chat messages parsed:', messages);
        } catch (parseError) {
            console.error('Failed to parse JSON response:', parseError);
            console.error('Parse error details:', {
                message: parseError.message,
                name: parseError.name
            });
            messagesContainer.innerHTML = '<div class="empty-state"><p>Ошибка при обработке ответа сервера</p></div>';
            return;
        }

        if (!Array.isArray(messages)) {
            console.error('Messages is not an array:', messages, 'Type:', typeof messages);
            messagesContainer.innerHTML = '<div class="empty-state"><p>Некорректный формат данных</p></div>';
            return;
        }

        console.log('Messages count:', messages.length);
        if (messages.length === 0) {
            messagesContainer.innerHTML = '<div class="empty-state"><p>Нет сообщений</p></div>';
        } else {
            // Очищаем контейнер
            messagesContainer.innerHTML = '';

            // Создаем все сообщения асинхронно
            const messagePromises = messages.map(async (message, index) => {
                try {
                    console.log(`Processing message ${index}:`, message);
                    const messageElement = await createChatMessage(message);
                    return messageElement;
                } catch (msgError) {
                    console.error(`Error creating message element ${index}:`, msgError, message);
                    // Создаем fallback элемент
                    const fallbackDiv = document.createElement('div');
                    fallbackDiv.className = 'chat-message';
                    fallbackDiv.innerHTML = `
                        <div class="chat-message-header">
                            <span class="chat-message-author">Unknown</span>
                            <span class="chat-message-time">N/A</span>
                        </div>
                        <div class="chat-message-text">${escapeHtml(message.text || message.Text || 'Ошибка загрузки сообщения')}</div>
                    `;
                    return fallbackDiv;
                }
            });

            // Ждем все промисы и добавляем элементы
            const messageElements = await Promise.all(messagePromises);
            messageElements.forEach(element => {
                if (element) {
                    messagesContainer.appendChild(element);
                }
            });

            // Прокручиваем вниз после добавления всех сообщений
            setTimeout(() => {
                messagesContainer.scrollTop = messagesContainer.scrollHeight;
            }, 100);
        }
    } catch (error) {
        console.error('Error loading chat messages:', error);
        console.error('Error details:', {
            message: error.message,
            stack: error.stack,
            name: error.name
        });

        let errorMessage = 'Ошибка соединения с сервером';
        if (error.message && error.message.includes('Failed to fetch')) {
            errorMessage = 'Не удалось подключиться к серверу. Убедитесь, что сервер запущен на http://localhost:8080';
        } else if (error.message) {
            errorMessage = 'Ошибка: ' + error.message;
        }

        messagesContainer.innerHTML = `<div class="empty-state"><p>${escapeHtml(errorMessage)}</p></div>`;
    }
}

async function createChatMessage(message) {
    // Поддерживаем оба варианта имен полей
    const userId = message.user_id || message.userID || message.UserID;
    const messageText = message.text || message.Text || '';
    const createdAt = message.created_at || message.CreatedAt || message.createdAt;

    if (!userId) {
        console.error('Message missing user_id:', message);
    }

    // Используем кэш для пользователей
    let userName = 'Unknown';
    if (userId) {
        if (userCache[userId]) {
            userName = userCache[userId];
        } else {
            try {
                const userResponse = await usersAPI.getUserById(userId);
                if (userResponse.ok) {
                    const user = await userResponse.json();
                    userName = user.full_name || user.FullName || user.email || user.Email || 'Unknown';
                    userCache[userId] = userName; // Кэшируем
                }
            } catch (error) {
                console.error('Error loading user:', error, 'user_id:', userId);
            }
        }
    }

    const messageDiv = document.createElement('div');
    messageDiv.className = 'chat-message';

    let formattedDate = 'N/A';
    if (createdAt) {
        try {
            formattedDate = new Date(createdAt).toLocaleString('ru-RU');
        } catch (e) {
            console.error('Error parsing date:', e, 'date:', createdAt);
        }
    }

    messageDiv.innerHTML = `
        <div class="chat-message-header">
            <span class="chat-message-author">${escapeHtml(userName)}</span>
            <span class="chat-message-time">${formattedDate}</span>
        </div>
        <div class="chat-message-text">${escapeHtml(messageText)}</div>
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

    if (!currentTaskId) {
        alert('Задача не выбрана');
        return;
    }

    try {
        console.log('Sending message:', { type: 'tasks', id: currentTaskId, userId, text });
        const response = await chatAPI.sendMessage('tasks', currentTaskId, userId, text);
        console.log('Message response status:', response.status, response.ok);

        if (!response.ok) {
            // Пытаемся получить текст ошибки
            let errorMessage = 'Ошибка при отправке сообщения';
            try {
                const errorText = await response.text();
                console.error('Error response:', errorText);
                let errorData = {};
                try {
                    errorData = JSON.parse(errorText);
                } catch (e) {
                    // Если не JSON, используем текст как есть
                    errorMessage = errorText || errorMessage;
                }
                if (errorData.error) {
                    errorMessage = errorData.error;
                }
            } catch (e) {
                console.error('Failed to read error response:', e);
            }

            if (response.status === 403) {
                // User doesn't have access
                const chatInputContainer = document.getElementById('chat-input-container');
                const chatNoAccess = document.getElementById('chat-no-access');
                if (chatInputContainer) chatInputContainer.style.display = 'none';
                if (chatNoAccess) chatNoAccess.style.display = 'block';
            } else {
                alert(errorMessage);
            }
            return;
        }

        // Успешная отправка
        const data = await response.json().catch(() => ({}));
        console.log('Message sent successfully:', data);

        const chatInput = document.getElementById('chat-input');
        if (chatInput) chatInput.value = '';

        // Перезагружаем сообщения сразу после отправки
        console.log('Reloading chat messages after send...');
        await loadChatMessages('tasks', currentTaskId);
    } catch (error) {
        console.error('Error sending message:', error);
        console.error('Error details:', {
            message: error.message,
            stack: error.stack,
            name: error.name
        });

        let errorMessage = 'Ошибка соединения с сервером';
        if (error.message && error.message.includes('Failed to fetch')) {
            errorMessage = 'Не удалось подключиться к серверу. Убедитесь, что сервер запущен на http://localhost:8080';
        } else if (error.message) {
            errorMessage = 'Ошибка: ' + error.message;
        }

        alert(errorMessage);
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

// Функция для получения текущих фильтров задач
function getTaskFilters() {
    const filters = {};

    const statusEl = document.getElementById('filter-status');
    const priorityEl = document.getElementById('filter-priority');
    const sortByEl = document.getElementById('filter-sort-by');
    const sortOrderEl = document.getElementById('filter-sort-order');
    const searchEl = document.getElementById('filter-search');

    if (statusEl && statusEl.value) {
        filters.status = statusEl.value;
    }

    if (priorityEl && priorityEl.value !== '') {
        filters.priority = parseInt(priorityEl.value);
    }

    if (sortByEl && sortByEl.value) {
        filters.sort_by = sortByEl.value;
    }

    if (sortOrderEl && sortOrderEl.value) {
        filters.sort_order = sortOrderEl.value;
    }

    if (searchEl && searchEl.value.trim()) {
        filters.search = searchEl.value.trim();
    }

    return filters;
}

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

    if (!currentTaskId) {
        console.error('Current task ID is missing');
        showError('edit-task-error', 'Задача не выбрана');
        return;
    }

    const titleInput = document.getElementById('edit-task-title-input');
    const descriptionInput = document.getElementById('edit-task-description-input');
    const statusInput = document.getElementById('edit-task-status-input');
    const priorityInput = document.getElementById('edit-task-priority-input');

    if (!titleInput || !statusInput) {
        console.error('Edit task form inputs not found');
        showError('edit-task-error', 'Ошибка: форма не найдена');
        return;
    }

    // Очищаем название от восклицательных знаков перед отправкой
    // Восклицательные знаки добавляются только при отображении в бэке
    let titleValue = titleInput.value.trim();
    // Убираем все восклицательные знаки из конца названия
    titleValue = titleValue.replace(/!+$/, '');

    const description = descriptionInput ? descriptionInput.value.trim() : '';
    const newStatus = statusInput.value;
    const priority = priorityInput ? (parseInt(priorityInput.value) || 0) : 0;

    if (!titleValue) {
        showError('edit-task-error', 'Название задачи обязательно');
        return;
    }

    // Проверяем, изменился ли статус. Если нет - не отправляем его
    let statusToSend = null;
    if (originalTaskStatus !== null && newStatus !== originalTaskStatus) {
        statusToSend = newStatus;
        console.log('Status changed:', { from: originalTaskStatus, to: newStatus });
    } else if (originalTaskStatus === null) {
        // Если исходный статус не был сохранен, отправляем новый (для обратной совместимости)
        statusToSend = newStatus;
    } else {
        console.log('Status unchanged, not sending status field');
    }

    try {
        console.log('Updating task:', { id: currentTaskId, title: titleValue, description, status: statusToSend, priority });
        const response = await tasksAPI.updateTask(currentTaskId, titleValue, description, statusToSend, priority);
        console.log('Update task response status:', response.status, response.ok);

        if (!response.ok) {
            let errorMessage = 'Ошибка при обновлении задачи';
            try {
                const errorText = await response.text();
                console.error('Failed to update task:', response.status, errorText);
                let errorData = {};
                try {
                    errorData = JSON.parse(errorText);
                } catch (e) {
                    console.error('Failed to parse error response');
                }
                if (errorData.error) {
                    errorMessage = errorData.error;
                    // Улучшаем сообщение об ошибке для пользователя
                    if (errorMessage.includes('status changes only in a certain order') ||
                        errorMessage.includes('the task status changes only in a certain order')) {
                        errorMessage = 'Недопустимый переход статуса. Правила: To Do → In Progress; In Progress → To Do или Done; Done → In Progress';
                    }
                } else if (errorText && errorText !== 'Unknown error') {
                    errorMessage = errorText;
                    if (errorMessage.includes('status changes only in a certain order')) {
                        errorMessage = 'Недопустимый переход статуса. Правила: To Do → In Progress; In Progress → To Do или Done; Done → In Progress';
                    }
                }
            } catch (e) {
                console.error('Failed to read error response:', e);
            }
            showError('edit-task-error', errorMessage);
            return;
        }

        const data = await response.json().catch(() => ({}));
        console.log('Task updated successfully:', data);

        // Обновляем исходный статус после успешного обновления
        if (statusToSend !== null) {
            originalTaskStatus = statusToSend;
        }

        closeEditTaskModal();
        showTask(currentTaskId); // Перезагружаем задачу
    } catch (error) {
        console.error('Error updating task:', error);
        console.error('Error details:', {
            message: error.message,
            stack: error.stack,
            name: error.name
        });
        showError('edit-task-error', 'Ошибка соединения с сервером: ' + error.message);
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
    if (!currentTaskId) {
        console.error('Cannot open edit modal: task ID is missing');
        alert('Ошибка: задача не выбрана');
        return;
    }

    // Заполняем форму текущими данными задачи
    const taskTitleEl = document.getElementById('task-title');
    const taskDescriptionEl = document.getElementById('task-description');
    const taskStatusEl = document.getElementById('task-status');
    const taskPriorityEl = document.getElementById('task-priority');

    if (!taskTitleEl || !taskStatusEl) {
        console.error('Task elements not found');
        alert('Ошибка: не удалось загрузить данные задачи');
        return;
    }

    const taskTitle = taskTitleEl.textContent.trim();
    const taskDescription = taskDescriptionEl ? taskDescriptionEl.textContent.trim() : '';
    const taskStatus = taskStatusEl.textContent.trim().toLowerCase();
    const taskPriorityText = taskPriorityEl ? taskPriorityEl.textContent.trim() : 'N/A';

    // Преобразуем статус
    let statusValue = taskStatus;
    if (taskStatus === 'in progress' || taskStatus === 'in-progress') {
        statusValue = 'in_progress';
    }

    // Сохраняем исходный статус для сравнения при сохранении
    originalTaskStatus = statusValue;

    // Преобразуем приоритет - сначала пытаемся взять из data-атрибута
    let priorityValue = 0;
    if (taskPriorityEl) {
        const dataPriority = taskPriorityEl.getAttribute('data-priority');
        if (dataPriority !== null && dataPriority !== '') {
            priorityValue = parseInt(dataPriority) || 0;
            console.log('Priority from data-attribute:', priorityValue);
        } else if (taskPriorityText && taskPriorityText !== 'N/A') {
            // Fallback: пытаемся извлечь число из текста
            if (taskPriorityText.includes('Очень важно')) {
                priorityValue = 2;
            } else if (taskPriorityText.includes('Важно')) {
                priorityValue = 1;
            } else {
                const priorityMatch = taskPriorityText.match(/\d+/);
                if (priorityMatch) {
                    priorityValue = parseInt(priorityMatch[0]);
                }
            }
            console.log('Priority from text:', priorityValue, 'text:', taskPriorityText);
        }
    }

    const titleInput = document.getElementById('edit-task-title-input');
    const descriptionInput = document.getElementById('edit-task-description-input');
    const statusInput = document.getElementById('edit-task-status-input');
    const priorityInput = document.getElementById('edit-task-priority-input');

    if (!titleInput || !statusInput) {
        console.error('Edit task form inputs not found');
        alert('Ошибка: форма редактирования не найдена');
        return;
    }

    titleInput.value = taskTitle;
    if (descriptionInput) descriptionInput.value = taskDescription;
    if (priorityInput) priorityInput.value = priorityValue;

    // Правила перехода статусов (соответствуют backend)
    const allowedTransitions = {
        'todo': ['in_progress'],
        'in_progress': ['todo', 'done'],
        'done': ['in_progress']
    };

    // Очищаем и заполняем select только допустимыми статусами
    statusInput.innerHTML = '';

    // Добавляем текущий статус как первый вариант
    const currentStatusOption = document.createElement('option');
    currentStatusOption.value = statusValue;
    currentStatusOption.textContent = statusValue === 'todo' ? 'To Do' :
        statusValue === 'in_progress' ? 'In Progress' : 'Done';
    currentStatusOption.selected = true;
    statusInput.appendChild(currentStatusOption);

    // Добавляем допустимые переходы
    const allowedStatuses = allowedTransitions[statusValue] || [];
    allowedStatuses.forEach(allowedStatus => {
        if (allowedStatus !== statusValue) {
            const option = document.createElement('option');
            option.value = allowedStatus;
            option.textContent = allowedStatus === 'todo' ? 'To Do' :
                allowedStatus === 'in_progress' ? 'In Progress' : 'Done';
            statusInput.appendChild(option);
        }
    });

    console.log('Edit task modal opened with data:', {
        title: taskTitle,
        description: taskDescription,
        status: statusValue,
        priority: priorityValue,
        allowedTransitions: allowedStatuses,
        originalStatus: originalTaskStatus
    });

    const editTaskModal = document.getElementById('edit-task-modal');
    if (editTaskModal) {
        editTaskModal.style.display = 'flex';
        titleInput.focus();
    } else {
        console.error('Edit task modal not found');
    }
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
    if (!reportsContent) {
        console.error('Reports content container not found');
        return;
    }

    if (!projectId) {
        console.error('Project ID is missing for reports');
        reportsContent.innerHTML = '<div class="empty-state"><p>ID проекта не указан</p></div>';
        return;
    }

    reportsContent.innerHTML = '<div class="loading">Загрузка отчетов...</div>';

    try {
        console.log('Loading reports for project ID:', projectId);
        const [topWorkersRes, avgTimeRes, completionRes] = await Promise.all([
            reportsAPI.getTopWorkers(projectId),
            reportsAPI.getAverageTime(projectId),
            reportsAPI.getCompletionPercent(projectId),
        ]);

        console.log('Reports responses:', {
            topWorkers: { status: topWorkersRes.status, ok: topWorkersRes.ok },
            avgTime: { status: avgTimeRes.status, ok: avgTimeRes.ok },
            completion: { status: completionRes.status, ok: completionRes.ok }
        });

        let html = '<div class="reports-grid">';

        // Top Workers
        if (topWorkersRes.ok) {
            try {
                const topWorkers = await topWorkersRes.json();
                console.log('Top workers data:', topWorkers);
                html += '<div class="report-card"><h4>Топ работников</h4>';
                if (Array.isArray(topWorkers) && topWorkers.length > 0) {
                    html += '<ul>';
                    topWorkers.forEach(worker => {
                        // Поддерживаем оба варианта имен полей
                        const name = worker.name || worker.Name || worker.full_name || worker.FullName || 'Unknown';
                        const tasks = worker.completed_tasks || worker.CompletedTasks || worker.task_count || worker.TaskCount || 0;
                        html += `<li>${escapeHtml(name)}: ${tasks} задач</li>`;
                    });
                    html += '</ul>';
                } else {
                    html += '<p>Нет данных</p>';
                }
                html += '</div>';
            } catch (e) {
                console.error('Error parsing top workers:', e);
                html += '<div class="report-card"><h4>Топ работников</h4><p>Ошибка загрузки</p></div>';
            }
        } else {
            console.error('Failed to load top workers:', topWorkersRes.status);
            const errorText = await topWorkersRes.text().catch(() => '');
            console.error('Error response:', errorText);
            html += '<div class="report-card"><h4>Топ работников</h4><p>Ошибка загрузки</p></div>';
        }

        // Average Time
        if (avgTimeRes.ok) {
            try {
                const avgTime = await avgTimeRes.json();
                console.log('Average time data:', avgTime);
                html += '<div class="report-card"><h4>Среднее время выполнения</h4>';
                // Поддерживаем оба варианта имен полей
                const avgHuman = avgTime.average_human || avgTime.AverageHuman || avgTime.average_time || avgTime.AverageTime;
                const tasksCount = avgTime.tasks_count || avgTime.TasksCount || avgTime.completed_count || avgTime.CompletedCount || 0;

                if (avgHuman) {
                    html += `<p>${escapeHtml(avgHuman)}</p>`;
                    if (tasksCount > 0) {
                        html += `<p class="report-detail">На основе ${tasksCount} задач</p>`;
                    }
                } else {
                    html += '<p>Нет данных</p>';
                }
                html += '</div>';
            } catch (e) {
                console.error('Error parsing average time:', e);
                html += '<div class="report-card"><h4>Среднее время выполнения</h4><p>Ошибка загрузки</p></div>';
            }
        } else {
            console.error('Failed to load average time:', avgTimeRes.status);
            html += '<div class="report-card"><h4>Среднее время выполнения</h4><p>Ошибка загрузки</p></div>';
        }

        // Completion Percent
        if (completionRes.ok) {
            try {
                const completion = await completionRes.json();
                console.log('Completion percent data:', completion);
                html += '<div class="report-card"><h4>Процент завершения</h4>';
                // Поддерживаем оба варианта имен полей
                const percent = completion.percent !== undefined ? completion.percent : completion.Percent;
                const totalTasks = completion.total_tasks || completion.TotalTasks || 0;
                const doneTasks = completion.done_tasks || completion.DoneTasks || 0;

                if (percent !== undefined) {
                    html += `<p class="report-percent">${percent.toFixed(1)}%</p>`;
                    html += `<p class="report-detail">${doneTasks} из ${totalTasks} задач</p>`;
                } else {
                    html += '<p>Нет данных</p>';
                }
                html += '</div>';
            } catch (e) {
                console.error('Error parsing completion percent:', e);
                html += '<div class="report-card"><h4>Процент завершения</h4><p>Ошибка загрузки</p></div>';
            }
        } else {
            console.error('Failed to load completion percent:', completionRes.status);
            html += '<div class="report-card"><h4>Процент завершения</h4><p>Ошибка загрузки</p></div>';
        }

        html += '</div>';
        reportsContent.innerHTML = html;
        console.log('Reports loaded successfully');
    } catch (error) {
        console.error('Error loading reports:', error);
        console.error('Error details:', {
            message: error.message,
            stack: error.stack,
            name: error.name
        });
        reportsContent.innerHTML = '<div class="empty-state"><p>Ошибка при загрузке отчетов: ' + escapeHtml(error.message) + '</p></div>';
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
    const editTaskForm = document.getElementById('edit-task-form');
    const closeEditTaskModalBtn = document.getElementById('close-edit-task-modal');
    const cancelEditTaskBtn = document.getElementById('cancel-edit-task');
    const editTaskModal = document.getElementById('edit-task-modal');

    if (editTaskBtn) {
        editTaskBtn.addEventListener('click', openEditTaskModal);
    }

    if (editTaskForm) {
        editTaskForm.addEventListener('submit', handleEditTask);
        console.log('Edit task form submit handler attached');
    } else {
        console.error('Edit task form not found!');
    }

    if (closeEditTaskModalBtn) {
        closeEditTaskModalBtn.addEventListener('click', closeEditTaskModal);
    }

    if (cancelEditTaskBtn) {
        cancelEditTaskBtn.addEventListener('click', closeEditTaskModal);
    }

    if (editTaskModal) {
        editTaskModal.addEventListener('click', (e) => {
            if (e.target.id === 'edit-task-modal') {
                closeEditTaskModal();
            }
        });
    }

    if (deleteTaskBtn) {
        deleteTaskBtn.addEventListener('click', handleDeleteTask);
    }

    // Task filters handlers
    const filterStatus = document.getElementById('filter-status');
    const filterPriority = document.getElementById('filter-priority');
    const filterSortBy = document.getElementById('filter-sort-by');
    const filterSortOrder = document.getElementById('filter-sort-order');
    const filterSearch = document.getElementById('filter-search');
    const clearFiltersBtn = document.getElementById('clear-filters');

    if (filterStatus) {
        filterStatus.addEventListener('change', () => {
            if (currentProjectId) {
                showProject(currentProjectId);
            }
        });
    }

    if (filterPriority) {
        filterPriority.addEventListener('change', () => {
            if (currentProjectId) {
                showProject(currentProjectId);
            }
        });
    }

    if (filterSortBy) {
        filterSortBy.addEventListener('change', () => {
            if (currentProjectId) {
                showProject(currentProjectId);
            }
        });
    }

    if (filterSortOrder) {
        filterSortOrder.addEventListener('change', () => {
            if (currentProjectId) {
                showProject(currentProjectId);
            }
        });
    }

    if (filterSearch) {
        let searchTimeout;
        filterSearch.addEventListener('input', () => {
            clearTimeout(searchTimeout);
            searchTimeout = setTimeout(() => {
                if (currentProjectId) {
                    showProject(currentProjectId);
                }
            }, 500); // Debounce search
        });
    }

    if (clearFiltersBtn) {
        clearFiltersBtn.addEventListener('click', () => {
            if (filterStatus) filterStatus.value = '';
            if (filterPriority) filterPriority.value = '';
            if (filterSortBy) filterSortBy.value = '';
            if (filterSortOrder) filterSortOrder.value = 'asc';
            if (filterSearch) filterSearch.value = '';
            if (currentProjectId) {
                showProject(currentProjectId);
            }
        });
    }
});

