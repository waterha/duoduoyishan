// API基础URL
const API_BASE_URL = '/api';

// 全局变量
let token = localStorage.getItem('token');
let currentUser = null;
let ws = null;
let currentChatUser = null;
let currentChatType = 'private'; // 'private' 或 'community'
let currentCommunity = null;

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', () => {
    if (token) {
        getUserInfo();
    }
    
    // 绑定表单提交事件
    document.getElementById('login-form').addEventListener('submit', handleLogin);
    document.getElementById('register-form').addEventListener('submit', handleRegister);
});

// 切换登录/注册标签
function switchTab(tab) {
    const tabs = document.querySelectorAll('.auth-tabs .tab-btn');
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    
    tabs.forEach(t => t.classList.remove('active'));
    
    if (tab === 'login') {
        tabs[0].classList.add('active');
        loginForm.style.display = 'flex';
        registerForm.style.display = 'none';
    } else {
        tabs[1].classList.add('active');
        loginForm.style.display = 'none';
        registerForm.style.display = 'flex';
    }
}

// 处理登录
async function handleLogin(e) {
    e.preventDefault();
    
    const username = document.getElementById('login-username').value;
    const password = document.getElementById('login-password').value;
    
    try {
        const response = await fetch(`${API_BASE_URL}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username, password })
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            token = data.data.token;
            currentUser = data.data.user;
            localStorage.setItem('token', token);
            showToast('登录成功', 'success');
            showMainPage();
            connectWebSocket();
        } else {
            showToast(data.message, 'error');
        }
    } catch (error) {
        showToast('登录失败，请检查网络连接', 'error');
    }
}

// 处理注册
async function handleRegister(e) {
    e.preventDefault();
    
    const username = document.getElementById('register-username').value;
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;
    
    try {
        const response = await fetch(`${API_BASE_URL}/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username, email, password })
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            showToast('注册成功，请登录', 'success');
            switchTab('login');
        } else {
            showToast(data.message, 'error');
        }
    } catch (error) {
        showToast('注册失败，请检查网络连接', 'error');
    }
}

// 获取用户信息
async function getUserInfo() {
    try {
        const response = await fetch(`${API_BASE_URL}/user/info`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            currentUser = data.data;
            showMainPage();
            connectWebSocket();
        } else {
            localStorage.removeItem('token');
            token = null;
        }
    } catch (error) {
        console.error('获取用户信息失败:', error);
    }
}

// 显示主页面
function showMainPage() {
    document.getElementById('auth-page').style.display = 'none';
    document.getElementById('main-page').style.display = 'flex';
    
    // 更新用户信息显示
    document.getElementById('user-nickname').textContent = currentUser.nickname || currentUser.username;
    
    // 加载数据
    loadFriends();
    loadCommunities();
    loadFriendRequests();
    loadPrivateChatList();
}

// 加载社区列表并自动选择第一个社区
async function loadCommunities() {
    try {
        const response = await fetch(`${API_BASE_URL}/community/my`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        const data = await response.json();
        
        if (data.code === 200 && data.data.communities && data.data.communities.length > 0) {
            const container = document.getElementById('community-list');
            container.innerHTML = '';
            
            data.data.communities.forEach(community => {
                const item = document.createElement('div');
                item.className = 'community-item';
                item.innerHTML = `
                    <div class="community-avatar">
                        <i class="fas fa-users"></i>
                    </div>
                    <div class="community-info">
                        <div class="community-name">${escapeHtml(community.name)}</div>
                        <div class="community-desc">${escapeHtml(community.description || '无描述')}</div>
                        <div class="community-meta">${community.member_count} 成员</div>
                    </div>
                    <button class="community-join-btn" onclick="joinCommunity(${community.id})">进入</button>
                `;
                container.appendChild(item);
            });
            
            // 自动选择第一个社区
            const firstCommunity = data.data.communities[0];
            if (firstCommunity) {
                startCommunityChat(firstCommunity.id, firstCommunity.name);
            }
        } else {
            document.getElementById('community-list').innerHTML = '<p style="text-align: center; color: #999; padding: 20px;">暂无社区<br><small>请先创建或加入社区</small></p>';
        }
    } catch (error) {
        console.error('加载社区列表失败:', error);
        document.getElementById('community-list').innerHTML = '<p style="text-align: center; color: #999; padding: 20px;">加载失败</p>';
    }
}

// 切换聊天类型
function switchChatType(type) {
    currentChatType = type;
    const buttons = document.querySelectorAll('.chat-type-btn');
    buttons.forEach(btn => btn.classList.remove('active'));
    event.currentTarget.classList.add('active');
    
    if (type === 'private') {
        loadPrivateChatList();
    } else {
        loadCommunityChatList();
    }
    
    // 重置聊天窗口
    document.getElementById('chat-window-title').textContent = '选择一个聊天';
    document.getElementById('chat-type-indicator').textContent = '';
    document.getElementById('chat-messages').innerHTML = `
        <div class="chat-empty">
            <i class="fas fa-comments"></i>
            <p>选择一个聊天开始对话</p>
        </div>
    `;
    currentChatUser = null;
    currentCommunity = null;
}

// 加载私聊列表
function loadPrivateChatList() {
    const container = document.getElementById('chat-list-items');
    container.innerHTML = '';
    container.setAttribute('data-type', 'private');
    
    // 从好友列表加载私聊
    loadFriends().then(() => {
        const friendsList = document.getElementById('friends-list');
        const friends = friendsList.querySelectorAll('.friend-card');
        
        if (friends.length === 0) {
            container.innerHTML = '<p style="text-align: center; color: #999; padding: 20px;">暂无聊天</p>';
            return;
        }
        
        friends.forEach(friend => {
            const friendName = friend.querySelector('.friend-name').textContent;
            const friendStatus = friend.querySelector('.friend-status').textContent;
            const chatBtn = friend.querySelector('.btn-small');
            const onclickAttr = chatBtn.getAttribute('onclick');
            const match = onclickAttr.match(/startChat\((\d+),/);
            const friendId = match ? match[1] : null;
            
            if (friendId) {
                const item = createChatItem(friendId, friendName, friendStatus, 'private');
                container.appendChild(item);
            }
        });
    });
}

// 加载社区聊天列表
async function loadCommunityChatList() {
    const container = document.getElementById('chat-list-items');
    container.innerHTML = '';
    container.setAttribute('data-type', 'community');
    
    try {
        const response = await fetch(`${API_BASE_URL}/community/my`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        const data = await response.json();
        
        if (data.code === 200 && data.data.communities && data.data.communities.length > 0) {
            data.data.communities.forEach(community => {
                const item = createChatItem(community.id, community.name, `${community.member_count} 成员`, 'community');
                container.appendChild(item);
            });
        } else {
            container.innerHTML = '<p style="text-align: center; color: #999; padding: 20px;">暂无社区聊天<br><small>请先加入社区</small></p>';
        }
    } catch (error) {
        console.error('加载社区聊天列表失败:', error);
        container.innerHTML = '<p style="text-align: center; color: #999; padding: 20px;">加载失败</p>';
    }
}

// 创建聊天项
function createChatItem(id, name, subtitle, type) {
    const item = document.createElement('div');
    item.className = 'chat-item';
    item.setAttribute('data-id', id);
    item.setAttribute('data-type', type);
    
    const icon = type === 'community' ? 'fa-users' : 'fa-user';
    
    item.innerHTML = `
        <div class="chat-item-avatar">
            <i class="fas ${icon}"></i>
        </div>
        <div class="chat-item-info">
            <div class="chat-item-name">${escapeHtml(name)}</div>
            <div class="chat-item-message">${escapeHtml(subtitle)}</div>
        </div>
    `;
    
    item.onclick = () => {
        if (type === 'community') {
            startCommunityChat(id, name);
        } else {
            startChat(id, name);
        }
    };
    
    return item;
}

// 退出登录
async function logout() {
    try {
        await fetch(`${API_BASE_URL}/auth/logout`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        localStorage.removeItem('token');
        token = null;
        currentUser = null;
        
        if (ws) {
            ws.close();
        }
        
        document.getElementById('auth-page').style.display = 'flex';
        document.getElementById('main-page').style.display = 'none';
        
        showToast('已退出登录', 'info');
    } catch (error) {
        showToast('退出失败', 'error');
    }
}

// 切换页面
function switchPage(page) {
    const pages = document.querySelectorAll('.page');
    const menuItems = document.querySelectorAll('.menu-item');
    
    pages.forEach(p => p.classList.remove('active'));
    menuItems.forEach(m => m.classList.remove('active'));
    
    document.getElementById(`${page}-page`).classList.add('active');
    event.currentTarget.classList.add('active');
}

// 连接WebSocket
function connectWebSocket() {
    if (!token) {
        console.log('没有token，无法连接WebSocket');
        return;
    }
    
    const wsUrl = `ws://${window.location.host}/api/ws?token=${token}`;
    ws = new WebSocket(wsUrl);
    
    ws.onopen = () => {
        console.log('WebSocket连接成功');
        showToast('连接成功', 'success');
    };
    
    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        handleWebSocketMessage(data);
    };
    
    ws.onerror = (error) => {
        console.error('WebSocket错误:', error);
        showToast('连接错误，请刷新页面重试', 'error');
    };
    
    ws.onclose = () => {
        console.log('WebSocket连接关闭');
        showToast('连接已断开', 'info');
        // 尝试重新连接
        setTimeout(() => {
            if (token) {
                connectWebSocket();
            }
        }, 3000);
    };
}

// 处理WebSocket消息
function handleWebSocketMessage(data) {
    switch (data.type) {
        case 'message':
            // 私聊消息
            if (data.data.to_type === 1 && currentChatUser && data.data.from_user_id === currentChatUser.id) {
                appendMessage(data.data, 'received');
            }
            // 社区消息
            else if (data.data.to_type === 2 && currentCommunity && data.data.to_id === currentCommunity.id) {
                const msgType = data.data.from_user_id === currentUser.id ? 'sent' : 'received';
                appendMessage(data.data, msgType);
            }
            break;
        case 'user_status':
            updateFriendStatus(data.user_id, data.online);
            break;
        case 'typing':
            // 显示正在输入提示
            break;
        case 'online_count':
            // 更新在线人数
            if (currentCommunity) {
                document.getElementById('chat-type-indicator').textContent = `社区 · ${data.count}人在线`;
            }
            break;
    }
}

// 发送消息
function sendMessage() {
    const input = document.getElementById('message-input');
    const content = input.value.trim();
    
    if (!content) {
        showToast('请输入消息内容', 'error');
        return;
    }
    
    if (!ws || ws.readyState !== WebSocket.OPEN) {
        showToast('WebSocket未连接，请刷新页面重试', 'error');
        return;
    }
    
    let message;
    
    if (currentChatType === 'community' && currentCommunity) {
        // 社区消息
        message = {
            type: 'message',
            data: {
                to_type: 2, // 社区聊天
                to_id: currentCommunity.id,
                msg_type: 1, // 文本消息
                content: content
            }
        };
    } else if (currentChatType === 'private' && currentChatUser) {
        // 私聊消息
        message = {
            type: 'message',
            data: {
                to_type: 1, // 私聊
                to_id: currentChatUser.id,
                msg_type: 1, // 文本消息
                content: content
            }
        };
    } else {
        showToast('请先选择聊天对象', 'error');
        return;
    }
    
    ws.send(JSON.stringify(message));
    
    // 显示发送的消息
    appendMessage({
        content: content,
        created_at: new Date().toISOString()
    }, 'sent');
    
    input.value = '';
}

// 处理键盘事件
function handleKeyPress(event) {
    if (event.key === 'Enter' && !event.shiftKey) {
        event.preventDefault();
        sendMessage();
    }
}

// 添加消息到聊天窗口
function appendMessage(message, type) {
    const messagesContainer = document.getElementById('chat-messages');
    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${type}`;
    
    const time = new Date(message.created_at).toLocaleTimeString('zh-CN', {
        hour: '2-digit',
        minute: '2-digit'
    });
    
    messageDiv.innerHTML = `
        <div class="message-content">${escapeHtml(message.content)}</div>
        <span class="message-time">${time}</span>
    `;
    
    messagesContainer.appendChild(messageDiv);
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
}

// HTML转义
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// 加载好友列表
async function loadFriends() {
    try {
        const response = await fetch(`${API_BASE_URL}/friend/list`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            renderFriendsList(data.data.friends);
            return data.data.friends;
        }
        return [];
    } catch (error) {
        console.error('加载好友列表失败:', error);
        return [];
    }
}

// 渲染好友列表
function renderFriendsList(friends) {
    const container = document.getElementById('friends-list');
    container.innerHTML = '';
    
    if (!friends || friends.length === 0) {
        container.innerHTML = '<p style="text-align: center; color: #999; padding: 20px;">暂无好友</p>';
        return;
    }
    
    friends.forEach(friend => {
        const card = document.createElement('div');
        card.className = 'friend-card';
        card.innerHTML = `
            <div class="friend-avatar">
                <i class="fas fa-user"></i>
            </div>
            <div class="friend-info">
                <div class="friend-name">${escapeHtml(friend.nickname || friend.username)}</div>
                <div class="friend-status">${friend.status === 1 ? '在线' : '离线'}</div>
            </div>
            <div class="friend-actions">
                <button class="btn-small btn-accept" onclick="startChat(${friend.id}, '${escapeHtml(friend.nickname || friend.username)}')">
                    <i class="fas fa-comment"></i> 聊天
                </button>
            </div>
        `;
        container.appendChild(card);
    });
}

// 加载好友请求
async function loadFriendRequests() {
    try {
        const response = await fetch(`${API_BASE_URL}/friend/requests`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            const badge = document.getElementById('friend-request-badge');
            if (data.data.requests && data.data.requests.length > 0) {
                badge.style.display = 'inline';
                badge.textContent = data.data.requests.length;
            } else {
                badge.style.display = 'none';
            }
            renderFriendRequests(data.data.requests);
        }
    } catch (error) {
        console.error('加载好友请求失败:', error);
    }
}

// 渲染好友请求
function renderFriendRequests(requests) {
    const container = document.getElementById('friend-requests');
    container.innerHTML = '';
    
    if (!requests || requests.length === 0) {
        container.innerHTML = '<p style="text-align: center; color: #999; padding: 20px;">暂无好友请求</p>';
        return;
    }
    
    requests.forEach(request => {
        const card = document.createElement('div');
        card.className = 'friend-card';
        card.innerHTML = `
            <div class="friend-avatar">
                <i class="fas fa-user"></i>
            </div>
            <div class="friend-info">
                <div class="friend-name">用户ID: ${request.from_user_id}</div>
                <div class="friend-status">${request.message || '请求添加好友'}</div>
            </div>
            <div class="friend-actions">
                <button class="btn-small btn-accept" onclick="handleFriendRequest(${request.id}, 1)">接受</button>
                <button class="btn-small btn-reject" onclick="handleFriendRequest(${request.id}, 2)">拒绝</button>
            </div>
        `;
        container.appendChild(card);
    });
}

// 处理好友请求
async function handleFriendRequest(requestId, status) {
    try {
        const response = await fetch(`${API_BASE_URL}/friend/request/${requestId}`, {
            method: 'PUT',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ status })
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            showToast(status === 1 ? '已添加好友' : '已拒绝请求', 'success');
            loadFriends();
            loadFriendRequests();
        } else {
            showToast(data.message, 'error');
        }
    } catch (error) {
        showToast('操作失败', 'error');
    }
}

// 切换好友标签
function switchFriendTab(tab) {
    const tabs = document.querySelectorAll('.friends-tabs .tab-btn');
    tabs.forEach(t => t.classList.remove('active'));
    event.currentTarget.classList.add('active');
    
    if (tab === 'list') {
        document.getElementById('friends-list').style.display = 'grid';
        document.getElementById('friend-requests').style.display = 'none';
    } else {
        document.getElementById('friends-list').style.display = 'none';
        document.getElementById('friend-requests').style.display = 'grid';
    }
}

// 显示添加好友模态框
function showAddFriend() {
    const modal = document.getElementById('modal');
    const modalTitle = document.getElementById('modal-title');
    const modalBody = document.getElementById('modal-body');
    
    modalTitle.textContent = '添加好友';
    modalBody.innerHTML = `
        <div class="form-group">
            <label>用户ID</label>
            <input type="number" id="add-friend-id" placeholder="输入用户ID">
        </div>
        <div class="form-group">
            <label>验证消息</label>
            <input type="text" id="add-friend-message" placeholder="验证消息（可选）">
        </div>
        <button class="btn-primary" onclick="sendFriendRequest()">发送请求</button>
    `;
    
    modal.classList.add('active');
}

// 发送好友请求
async function sendFriendRequest() {
    const toUserId = document.getElementById('add-friend-id').value;
    const message = document.getElementById('add-friend-message').value;
    
    if (!toUserId) {
        showToast('请输入用户ID', 'error');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE_URL}/friend/request`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                to_user_id: parseInt(toUserId),
                message: message
            })
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            showToast('好友请求已发送', 'success');
            closeModal();
        } else {
            showToast(data.message, 'error');
        }
    } catch (error) {
        showToast('发送失败', 'error');
    }
}

// 加载社区列表
async function loadCommunities() {
    try {
        const response = await fetch(`${API_BASE_URL}/community/list`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            renderCommunityList(data.data.list);
        }
    } catch (error) {
        console.error('加载社区列表失败:', error);
    }
}

// 渲染社区列表
function renderCommunityList(communities) {
    const container = document.getElementById('community-list');
    container.innerHTML = '';
    
    if (!communities || communities.length === 0) {
        container.innerHTML = '<p style="text-align: center; color: #999; padding: 20px;">暂无社区</p>';
        return;
    }
    
    communities.forEach(community => {
        const card = document.createElement('div');
        card.className = 'community-card';
        card.innerHTML = `
            <div class="community-cover">
                <i class="fas fa-users"></i>
            </div>
            <div class="community-content">
                <div class="community-name">${escapeHtml(community.name)}</div>
                <div class="community-desc">${escapeHtml(community.description || '暂无描述')}</div>
                <div class="community-meta">
                    <span class="community-members">
                        <i class="fas fa-user"></i> ${community.member_count} 成员
                    </span>
                    <button class="btn-small btn-accept" onclick="joinCommunity(${community.id})">加入</button>
                </div>
            </div>
        `;
        container.appendChild(card);
    });
}

// 切换社区标签
function switchCommunityTab(tab) {
    const tabs = document.querySelectorAll('.community-tabs .tab-btn');
    tabs.forEach(t => t.classList.remove('active'));
    event.currentTarget.classList.add('active');
    
    if (tab === 'all') {
        loadCommunities();
    } else {
        loadMyCommunities();
    }
}

// 加载我的社区
async function loadMyCommunities() {
    try {
        const response = await fetch(`${API_BASE_URL}/community/my`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            renderCommunityList(data.data.communities);
        }
    } catch (error) {
        console.error('加载我的社区失败:', error);
    }
}

// 显示创建社区模态框
function showCreateCommunity() {
    const modal = document.getElementById('modal');
    const modalTitle = document.getElementById('modal-title');
    const modalBody = document.getElementById('modal-body');
    
    modalTitle.textContent = '创建社区';
    modalBody.innerHTML = `
        <div class="form-group">
            <label>社区名称</label>
            <input type="text" id="community-name" placeholder="社区名称">
        </div>
        <div class="form-group">
            <label>社区描述</label>
            <textarea id="community-desc" placeholder="社区描述"></textarea>
        </div>
        <div class="form-group">
            <label>社区分类</label>
            <input type="text" id="community-category" placeholder="社区分类">
        </div>
        <button class="btn-primary" onclick="createCommunity()">创建社区</button>
    `;
    
    modal.classList.add('active');
}

// 创建社区
async function createCommunity() {
    const name = document.getElementById('community-name').value;
    const description = document.getElementById('community-desc').value;
    const category = document.getElementById('community-category').value;
    
    if (!name || !category) {
        showToast('请填写社区名称和分类', 'error');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE_URL}/community/create`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ name, description, category })
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            showToast('社区创建成功', 'success');
            closeModal();
            loadCommunities();
        } else {
            showToast(data.message, 'error');
        }
    } catch (error) {
        showToast('创建失败', 'error');
    }
}

// 加入社区
async function joinCommunity(communityId) {
    try {
        const response = await fetch(`${API_BASE_URL}/community/join/${communityId}`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            showToast('加入成功', 'success');
            loadCommunities();
        } else {
            showToast(data.message, 'error');
        }
    } catch (error) {
        showToast('加入失败', 'error');
    }
}

// 开始聊天
function startChat(userId, username) {
    currentChatUser = { id: userId, username: username };
    currentCommunity = null;
    currentChatType = 'private';
    
    // 切换到聊天页面
    const menuItems = document.querySelectorAll('.menu-item');
    menuItems.forEach(m => m.classList.remove('active'));
    menuItems[0].classList.add('active');
    
    document.querySelectorAll('.page').forEach(p => p.classList.remove('active'));
    document.getElementById('chat-page').classList.add('active');
    
    // 更新聊天窗口标题
    document.getElementById('chat-window-title').textContent = username;
    document.getElementById('chat-type-indicator').textContent = '私聊';
    
    // 清空消息并加载历史消息
    document.getElementById('chat-messages').innerHTML = '';
    loadChatHistory(userId, 'private');
    
    // 高亮当前聊天项
    highlightChatItem(userId, 'private');
}

// 开始社区聊天
function startCommunityChat(communityId, communityName) {
    currentCommunity = { id: communityId, name: communityName };
    currentChatUser = null;
    currentChatType = 'community';
    
    // 切换到聊天页面
    const menuItems = document.querySelectorAll('.menu-item');
    menuItems.forEach(m => m.classList.remove('active'));
    menuItems[0].classList.add('active');
    
    document.querySelectorAll('.page').forEach(p => p.classList.remove('active'));
    document.getElementById('chat-page').classList.add('active');
    
    // 更新聊天窗口标题
    document.getElementById('chat-window-title').textContent = communityName;
    document.getElementById('chat-type-indicator').textContent = '社区';
    
    // 清空消息并加载历史消息
    document.getElementById('chat-messages').innerHTML = '';
    loadChatHistory(communityId, 'community');
    
    // 高亮当前聊天项
    highlightChatItem(communityId, 'community');
    
    // 获取在线人数
    getOnlineCount(communityId);
}

// 加载聊天历史
async function loadChatHistory(id, type) {
    try {
        const response = await fetch(`${API_BASE_URL}/message/history?to_type=${type === 'community' ? 2 : 1}&to_id=${id}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        const data = await response.json();
        
        if (data.code === 200 && data.data.messages) {
            data.data.messages.forEach(msg => {
                const msgType = msg.from_user_id === currentUser.id ? 'sent' : 'received';
                appendMessage(msg, msgType);
            });
        }
    } catch (error) {
        console.error('加载聊天历史失败:', error);
    }
}

// 高亮聊天项
function highlightChatItem(id, type) {
    const items = document.querySelectorAll('.chat-item');
    items.forEach(item => {
        item.classList.remove('active');
        if (item.getAttribute('data-id') == id && item.getAttribute('data-type') === type) {
            item.classList.add('active');
        }
    });
}

// 获取在线人数
async function getOnlineCount(communityId) {
    try {
        const response = await fetch(`${API_BASE_URL}/ws/online?room_id=group:${communityId}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            document.getElementById('chat-type-indicator').textContent = `社区 · ${data.data.online_count}人在线`;
        }
    } catch (error) {
        console.error('获取在线人数失败:', error);
    }
}

// 添加聊天项
function addChatItem(userId, username) {
    const container = document.getElementById('chat-list-items');
    
    // 检查是否已存在
    const existingItem = container.querySelector(`[data-user-id="${userId}"]`);
    if (existingItem) return;
    
    const item = document.createElement('div');
    item.className = 'chat-item';
    item.setAttribute('data-user-id', userId);
    item.onclick = () => startChat(userId, username);
    item.innerHTML = `
        <div class="chat-item-avatar">
            <i class="fas fa-user"></i>
        </div>
        <div class="chat-item-info">
            <div class="chat-item-name">${escapeHtml(username)}</div>
            <div class="chat-item-message">点击开始聊天</div>
        </div>
    `;
    
    container.insertBefore(item, container.firstChild);
}

// 更新好友状态
function updateFriendStatus(userId, online) {
    const statusElements = document.querySelectorAll('.friend-card');
    statusElements.forEach(el => {
        const statusEl = el.querySelector('.friend-status');
        if (statusEl) {
            statusEl.textContent = online ? '在线' : '离线';
        }
    });
}

// 更新个人资料
async function updateProfile() {
    const nickname = document.getElementById('profile-nickname').value;
    const signature = document.getElementById('profile-signature').value;
    const gender = document.getElementById('profile-gender').value;
    
    try {
        const response = await fetch(`${API_BASE_URL}/user/info`, {
            method: 'PUT',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ nickname, signature, gender: parseInt(gender) })
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            showToast('资料更新成功', 'success');
            currentUser = data.data;
            document.getElementById('user-nickname').textContent = currentUser.nickname || currentUser.username;
        } else {
            showToast(data.message, 'error');
        }
    } catch (error) {
        showToast('更新失败', 'error');
    }
}

// 修改密码
async function changePassword() {
    const oldPassword = document.getElementById('old-password').value;
    const newPassword = document.getElementById('new-password').value;
    
    if (!oldPassword || !newPassword) {
        showToast('请填写完整信息', 'error');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE_URL}/user/password`, {
            method: 'PUT',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ old_password: oldPassword, new_password: newPassword })
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            showToast('密码修改成功，请重新登录', 'success');
            setTimeout(() => logout(), 2000);
        } else {
            showToast(data.message, 'error');
        }
    } catch (error) {
        showToast('修改失败', 'error');
    }
}

// 关闭模态框
function closeModal() {
    document.getElementById('modal').classList.remove('active');
}

// 显示提示消息
function showToast(message, type = 'info') {
    const toast = document.getElementById('toast');
    toast.textContent = message;
    toast.className = `toast ${type} show`;
    
    setTimeout(() => {
        toast.classList.remove('show');
    }, 3000);
}

// 点击模态框外部关闭
document.getElementById('modal').addEventListener('click', (e) => {
    if (e.target === document.getElementById('modal')) {
        closeModal();
    }
});
