/**
 * 聊天 WebSocket JavaScript
 * 包含WebSocket連接、訊息處理、即時更新、連接管理
 */

class ChatWebSocketManager {
    constructor() {
        this.websocket = null;
        this.isConnected = false;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectDelay = 1000; // 初始重連延遲
        this.heartbeatInterval = null;
        this.messageQueue = []; // 離線訊息隊列
        this.currentChatId = null;
        this.token = localStorage.getItem('token');
        this.eventHandlers = new Map();
        
        this.init();
    }

    async init() {
        if (!this.token) {
            console.warn('未找到認證令牌，無法建立WebSocket連接');
            return;
        }

        this.bindUIEvents();
        this.connectWebSocket();
        this.setupConnectionMonitoring();
    }

    bindUIEvents() {
        // 發送訊息按鈕
        const sendBtn = document.getElementById('send-btn');
        if (sendBtn) {
            sendBtn.addEventListener('click', () => {
                this.sendMessage();
            });
        }

        // 訊息輸入框
        const messageInput = document.getElementById('message-input');
        if (messageInput) {
            messageInput.addEventListener('keydown', (e) => {
                if (e.key === 'Enter' && !e.shiftKey) {
                    e.preventDefault();
                    this.sendMessage();
                } else if (e.key !== 'Enter') {
                    // 發送正在輸入指示
                    this.sendTypingIndicator();
                }
            });

            // 停止輸入指示器
            let typingTimer;
            messageInput.addEventListener('input', () => {
                clearTimeout(typingTimer);
                typingTimer = setTimeout(() => {
                    this.stopTypingIndicator();
                }, 1000);
            });
        }

        // 連接狀態顯示
        this.updateConnectionStatus(false);
    }

    connectWebSocket() {
        try {
            // 建立WebSocket連接
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = `${protocol}//${window.location.host}/ws?token=${encodeURIComponent(this.token)}`;
            
            this.websocket = new WebSocket(wsUrl);
            
            this.websocket.onopen = this.onWebSocketOpen.bind(this);
            this.websocket.onmessage = this.onWebSocketMessage.bind(this);
            this.websocket.onclose = this.onWebSocketClose.bind(this);
            this.websocket.onerror = this.onWebSocketError.bind(this);
            
        } catch (error) {
            console.error('WebSocket連接失敗:', error);
            this.scheduleReconnect();
        }
    }

    onWebSocketOpen(event) {
        console.log('WebSocket連接已建立');
        this.isConnected = true;
        this.reconnectAttempts = 0;
        this.reconnectDelay = 1000;
        
        this.updateConnectionStatus(true);
        this.startHeartbeat();
        
        // 發送待發送的訊息
        this.flushMessageQueue();
        
        // 觸發連接成功事件
        this.emitEvent('connected');
    }

    onWebSocketMessage(event) {
        try {
            const data = JSON.parse(event.data);
            this.handleWebSocketMessage(data);
        } catch (error) {
            console.error('解析WebSocket訊息失敗:', error);
        }
    }

    onWebSocketClose(event) {
        console.log('WebSocket連接已關閉', event.code, event.reason);
        this.isConnected = false;
        this.updateConnectionStatus(false);
        this.stopHeartbeat();
        
        if (event.code !== 1000) { // 非正常關閉
            this.scheduleReconnect();
        }
        
        this.emitEvent('disconnected', event);
    }

    onWebSocketError(event) {
        console.error('WebSocket錯誤:', event);
        this.emitEvent('error', event);
    }

    handleWebSocketMessage(data) {
        switch (data.type) {
            case 'message':
                this.handleNewMessage(data.data);
                break;
            case 'typing_start':
                this.handleTypingStart(data.data);
                break;
            case 'typing_stop':
                this.handleTypingStop(data.data);
                break;
            case 'match_notification':
                this.handleMatchNotification(data.data);
                break;
            case 'user_online':
                this.handleUserOnline(data.data);
                break;
            case 'user_offline':
                this.handleUserOffline(data.data);
                break;
            case 'message_read':
                this.handleMessageRead(data.data);
                break;
            case 'pong':
                // 心跳響應，無需處理
                break;
            default:
                console.log('收到未知類型的訊息:', data);
        }
    }

    handleNewMessage(messageData) {
        // 檢查是否為當前聊天的訊息
        if (messageData.chatId === this.currentChatId) {
            this.addMessageToUI(messageData);
            this.scrollToBottom();
            
            // 發送已讀確認
            this.sendMessageRead(messageData.messageId);
        }
        
        // 更新聊天列表中的最後訊息
        this.updateChatListLastMessage(messageData);
        
        // 播放通知音效
        this.playNotificationSound();
        
        // 觸發新訊息事件
        this.emitEvent('newMessage', messageData);
    }

    handleTypingStart(data) {
        if (data.chatId === this.currentChatId) {
            this.showTypingIndicator(data.userName);
        }
    }

    handleTypingStop(data) {
        if (data.chatId === this.currentChatId) {
            this.hideTypingIndicator();
        }
    }

    handleMatchNotification(matchData) {
        this.showMatchNotification(matchData);
        this.emitEvent('newMatch', matchData);
    }

    handleUserOnline(data) {
        this.updateUserOnlineStatus(data.userId, true);
        this.emitEvent('userOnline', data);
    }

    handleUserOffline(data) {
        this.updateUserOnlineStatus(data.userId, false);
        this.emitEvent('userOffline', data);
    }

    handleMessageRead(data) {
        this.updateMessageReadStatus(data.messageId);
        this.emitEvent('messageRead', data);
    }

    sendMessage() {
        const messageInput = document.getElementById('message-input');
        if (!messageInput) return;
        
        const content = messageInput.value.trim();
        if (!content || !this.currentChatId) return;

        const messageData = {
            type: 'send_message',
            data: {
                chatId: this.currentChatId,
                content: content,
                messageType: 'text',
                timestamp: Date.now()
            }
        };

        if (this.isConnected) {
            this.websocket.send(JSON.stringify(messageData));
            messageInput.value = '';
            this.autoResizeTextarea(messageInput);
        } else {
            // 離線時加入隊列
            this.messageQueue.push(messageData);
            this.showOfflineNotice();
        }

        // 樂觀更新UI
        this.addMessageToUI({
            id: 'temp_' + Date.now(),
            content: content,
            senderId: this.getCurrentUserId(),
            timestamp: new Date().toISOString(),
            status: 'sending'
        });
    }

    sendTypingIndicator() {
        if (!this.isConnected || !this.currentChatId) return;

        const data = {
            type: 'typing_start',
            data: {
                chatId: this.currentChatId
            }
        };

        this.websocket.send(JSON.stringify(data));
    }

    stopTypingIndicator() {
        if (!this.isConnected || !this.currentChatId) return;

        const data = {
            type: 'typing_stop',
            data: {
                chatId: this.currentChatId
            }
        };

        this.websocket.send(JSON.stringify(data));
    }

    sendMessageRead(messageId) {
        if (!this.isConnected) return;

        const data = {
            type: 'message_read',
            data: {
                messageId: messageId
            }
        };

        this.websocket.send(JSON.stringify(data));
    }

    addMessageToUI(messageData) {
        const messagesContainer = document.getElementById('chat-messages');
        if (!messagesContainer) return;

        const isOwn = messageData.senderId === this.getCurrentUserId();
        const messageElement = document.createElement('div');
        messageElement.className = `message ${isOwn ? 'sent' : 'received'}`;
        messageElement.dataset.messageId = messageData.id;

        const time = this.formatTime(messageData.timestamp);
        const status = this.formatMessageStatus(messageData.status);

        messageElement.innerHTML = `
            <div class="message-content">
                ${this.escapeHtml(messageData.content)}
                <div class="message-time">${time}</div>
            </div>
            ${isOwn ? `<div class="message-status">${status}</div>` : ''}
        `;

        messagesContainer.appendChild(messageElement);
        this.scrollToBottom();

        // 如果是臨時訊息，添加發送動畫
        if (messageData.id && messageData.id.startsWith('temp_')) {
            messageElement.classList.add('sending');
        }
    }

    updateMessageReadStatus(messageId) {
        const messageElement = document.querySelector(`[data-message-id="${messageId}"]`);
        if (!messageElement) return;

        const statusElement = messageElement.querySelector('.message-status');
        if (statusElement) {
            statusElement.innerHTML = '✓✓';
            statusElement.classList.add('read');
        }
    }

    showTypingIndicator(userName) {
        const indicator = document.getElementById('typing-indicator');
        if (!indicator) return;

        const userSpan = indicator.querySelector('#typing-user');
        if (userSpan) {
            userSpan.textContent = userName || '對方';
        }

        indicator.classList.remove('hidden');
    }

    hideTypingIndicator() {
        const indicator = document.getElementById('typing-indicator');
        if (indicator) {
            indicator.classList.add('hidden');
        }
    }

    updateChatListLastMessage(messageData) {
        // 更新聊天列表中對應聊天的最後訊息
        const chatItem = document.querySelector(`[data-chat-id="${messageData.chatId}"]`);
        if (!chatItem) return;

        const lastMessageElement = chatItem.querySelector('.match-last-message');
        if (lastMessageElement) {
            lastMessageElement.textContent = messageData.content;
        }

        const timeElement = chatItem.querySelector('.match-time');
        if (timeElement) {
            timeElement.textContent = this.formatTime(messageData.timestamp);
        }

        // 如果不是當前聊天，增加未讀計數
        if (messageData.chatId !== this.currentChatId) {
            this.incrementUnreadCount(messageData.chatId);
        }
    }

    incrementUnreadCount(chatId) {
        const chatItem = document.querySelector(`[data-chat-id="${chatId}"]`);
        if (!chatItem) return;

        const badge = chatItem.querySelector('.unread-badge');
        if (badge) {
            const current = parseInt(badge.textContent) || 0;
            badge.textContent = current + 1;
            badge.classList.add('show');
        }
    }

    updateUserOnlineStatus(userId, isOnline) {
        // 更新聊天列表中用戶的在線狀態
        const userElements = document.querySelectorAll(`[data-user-id="${userId}"]`);
        userElements.forEach(element => {
            const onlineIndicator = element.querySelector('.online-indicator');
            if (onlineIndicator) {
                onlineIndicator.classList.toggle('show', isOnline);
            }
        });

        // 更新聊天標題中的狀態
        if (this.currentChatId) {
            const statusElement = document.getElementById('chat-user-status');
            if (statusElement && this.isCurrentChatUser(userId)) {
                statusElement.textContent = isOnline ? '在線' : '離線';
            }
        }
    }

    startHeartbeat() {
        this.stopHeartbeat(); // 確保不會重複建立

        this.heartbeatInterval = setInterval(() => {
            if (this.isConnected && this.websocket.readyState === WebSocket.OPEN) {
                this.websocket.send(JSON.stringify({ type: 'ping' }));
            }
        }, 30000); // 每30秒發送一次心跳
    }

    stopHeartbeat() {
        if (this.heartbeatInterval) {
            clearInterval(this.heartbeatInterval);
            this.heartbeatInterval = null;
        }
    }

    scheduleReconnect() {
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            console.error('WebSocket重連次數已達上限，停止重連');
            this.showPermanentDisconnectNotice();
            return;
        }

        this.reconnectAttempts++;
        const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1); // 指數退避

        console.log(`${delay}ms後嘗試第${this.reconnectAttempts}次重連...`);
        this.showReconnectingNotice(this.reconnectAttempts);

        setTimeout(() => {
            if (!this.isConnected) {
                this.connectWebSocket();
            }
        }, delay);
    }

    flushMessageQueue() {
        while (this.messageQueue.length > 0) {
            const messageData = this.messageQueue.shift();
            this.websocket.send(JSON.stringify(messageData));
        }
    }

    setupConnectionMonitoring() {
        // 監控網路狀態
        window.addEventListener('online', () => {
            console.log('網路已恢復，嘗試重新連接');
            if (!this.isConnected) {
                this.reconnectAttempts = 0; // 重置重連次數
                this.connectWebSocket();
            }
        });

        window.addEventListener('offline', () => {
            console.log('網路已斷開');
            this.showOfflineNotice();
        });

        // 頁面可見性變更
        document.addEventListener('visibilitychange', () => {
            if (!document.hidden && !this.isConnected) {
                // 頁面變為可見且未連接時嘗試重連
                this.connectWebSocket();
            }
        });
    }

    updateConnectionStatus(connected) {
        const statusElement = document.getElementById('connection-status');
        const statusText = document.getElementById('connection-text');
        
        if (statusElement) {
            statusElement.className = `connection-status ${connected ? 'connected' : 'disconnected'}`;
        }
        
        if (statusText) {
            statusText.textContent = connected ? '已連接' : '連接中斷';
        }

        // 更新發送按鈕狀態
        const sendBtn = document.getElementById('send-btn');
        if (sendBtn) {
            sendBtn.disabled = !connected;
        }
    }

    showReconnectingNotice(attempt) {
        this.showNotice(`正在重新連接... (${attempt}/${this.maxReconnectAttempts})`, 'warning');
    }

    showOfflineNotice() {
        this.showNotice('您目前處於離線狀態，訊息將在連接恢復後發送', 'warning', true);
    }

    showPermanentDisconnectNotice() {
        this.showNotice('連接失敗，請重新整理頁面', 'error', true);
    }

    showNotice(message, type = 'info', persistent = false) {
        // 移除現有通知
        const existingNotice = document.querySelector('.connection-notice');
        if (existingNotice) {
            existingNotice.remove();
        }

        const notice = document.createElement('div');
        notice.className = `connection-notice notice-${type}`;
        notice.textContent = message;

        document.body.appendChild(notice);

        if (!persistent) {
            setTimeout(() => {
                if (notice.parentNode) {
                    notice.remove();
                }
            }, 5000);
        }
    }

    showMatchNotification(matchData) {
        // 這個功能可以與配對頁面的通知整合
        console.log('收到配對通知:', matchData);
        
        // 可以在這裡顯示一個簡單的通知
        this.showNotice(`🎉 與 ${matchData.otherUser.name} 配對成功！`, 'success');
    }

    playNotificationSound() {
        // 播放通知音效（如果用戶允許）
        if ('Notification' in window && Notification.permission === 'granted') {
            // 可以播放音效或顯示桌面通知
            try {
                const audio = new Audio('/static/audio/notification.mp3');
                audio.volume = 0.3;
                audio.play().catch(() => {
                    // 忽略播放失敗
                });
            } catch (error) {
                // 音效文件不存在或播放失敗
            }
        }
    }

    scrollToBottom() {
        const messagesContainer = document.getElementById('chat-messages');
        if (messagesContainer) {
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
        }
    }

    autoResizeTextarea(textarea) {
        textarea.style.height = 'auto';
        textarea.style.height = Math.min(textarea.scrollHeight, 100) + 'px';
    }

    // 工具方法
    getCurrentUserId() {
        const user = JSON.parse(localStorage.getItem('user') || '{}');
        return user.id || null;
    }

    isCurrentChatUser(userId) {
        // 檢查是否為當前聊天的對方用戶
        // 這需要根據實際的聊天數據結構來實作
        return false; // 暫時返回false
    }

    formatTime(timestamp) {
        const date = new Date(timestamp);
        const now = new Date();
        
        if (now.toDateString() === date.toDateString()) {
            return date.toLocaleTimeString('zh-TW', { hour: '2-digit', minute: '2-digit' });
        } else {
            return date.toLocaleDateString('zh-TW', { month: 'short', day: 'numeric' });
        }
    }

    formatMessageStatus(status) {
        switch (status) {
            case 'sending':
                return '⏳';
            case 'sent':
                return '✓';
            case 'delivered':
                return '✓';
            case 'read':
                return '✓✓';
            default:
                return '';
        }
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    // 事件系統
    on(event, handler) {
        if (!this.eventHandlers.has(event)) {
            this.eventHandlers.set(event, []);
        }
        this.eventHandlers.get(event).push(handler);
    }

    off(event, handler) {
        if (!this.eventHandlers.has(event)) return;
        
        const handlers = this.eventHandlers.get(event);
        const index = handlers.indexOf(handler);
        if (index > -1) {
            handlers.splice(index, 1);
        }
    }

    emitEvent(event, data) {
        if (!this.eventHandlers.has(event)) return;
        
        this.eventHandlers.get(event).forEach(handler => {
            try {
                handler(data);
            } catch (error) {
                console.error(`事件處理器錯誤 (${event}):`, error);
            }
        });
    }

    // 公共API方法
    setCurrentChat(chatId) {
        this.currentChatId = chatId;
        
        // 通知服務器用戶進入了特定聊天
        if (this.isConnected && chatId) {
            const data = {
                type: 'join_chat',
                data: { chatId: chatId }
            };
            this.websocket.send(JSON.stringify(data));
        }
    }

    disconnect() {
        this.stopHeartbeat();
        if (this.websocket) {
            this.websocket.close(1000, '用戶主動斷開');
        }
    }

    reconnect() {
        this.disconnect();
        this.reconnectAttempts = 0;
        this.connectWebSocket();
    }
}

// 全域實例
let chatWebSocket;

// 初始化
document.addEventListener('DOMContentLoaded', () => {
    chatWebSocket = new ChatWebSocketManager();
    
    // 將實例綁定到全域，供其他頁面組件使用
    window.chatWebSocket = chatWebSocket;
});

// 頁面卸載時清理
window.addEventListener('beforeunload', () => {
    if (chatWebSocket) {
        chatWebSocket.disconnect();
    }
});

// 匯出供其他模組使用
if (typeof module !== 'undefined' && module.exports) {
    module.exports = ChatWebSocketManager;
}