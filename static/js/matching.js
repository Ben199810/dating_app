/**
 * 配對介面 JavaScript
 * 包含滑動邏輯、配對API、動畫效果、觸控支援
 */

class MatchingInterface {
    constructor() {
        this.cards = [];
        this.currentCardIndex = 0;
        this.isDragging = false;
        this.startX = 0;
        this.startY = 0;
        this.currentX = 0;
        this.currentY = 0;
        this.swipeThreshold = 100; // 滑動閾值
        this.rotationFactor = 0.1; // 旋轉係數
        this.isLoading = false;
        this.matchQueue = [];
        
        this.init();
    }

    async init() {
        try {
            this.bindEvents();
            this.setupGestures();
            await this.loadPotentialMatches();
            this.showCurrentCard();
        } catch (error) {
            console.error('初始化失敗:', error);
            this.showError('載入配對頁面失敗');
        }
    }

    bindEvents() {
        // 按鈕事件
        document.getElementById('pass-btn')?.addEventListener('click', () => {
            this.swipeLeft();
        });

        document.getElementById('like-btn')?.addEventListener('click', () => {
            this.swipeRight();
        });

        document.getElementById('super-like-btn')?.addEventListener('click', () => {
            this.superLike();
        });

        document.getElementById('rewind-btn')?.addEventListener('click', () => {
            this.rewind();
        });

        // 鍵盤支援
        document.addEventListener('keydown', (e) => {
            this.handleKeyPress(e);
        });

        // 載入更多按鈕
        document.getElementById('load-more')?.addEventListener('click', () => {
            this.loadPotentialMatches();
        });

        // 篩選器
        document.getElementById('filter-btn')?.addEventListener('click', () => {
            this.toggleFilters();
        });

        document.getElementById('apply-filters')?.addEventListener('click', () => {
            this.applyFilters();
        });
    }

    setupGestures() {
        const cardContainer = document.getElementById('cards-container');
        if (!cardContainer) return;

        // 觸控事件
        cardContainer.addEventListener('touchstart', this.handleTouchStart.bind(this), { passive: false });
        cardContainer.addEventListener('touchmove', this.handleTouchMove.bind(this), { passive: false });
        cardContainer.addEventListener('touchend', this.handleTouchEnd.bind(this), { passive: false });

        // 滑鼠事件（桌面支援）
        cardContainer.addEventListener('mousedown', this.handleMouseDown.bind(this));
        cardContainer.addEventListener('mousemove', this.handleMouseMove.bind(this));
        cardContainer.addEventListener('mouseup', this.handleMouseUp.bind(this));
        cardContainer.addEventListener('mouseleave', this.handleMouseUp.bind(this));
    }

    async loadPotentialMatches() {
        if (this.isLoading) return;
        
        this.isLoading = true;
        this.showLoading(true);

        try {
            const token = localStorage.getItem('token');
            if (!token) {
                window.location.href = '/static/html/register.html';
                return;
            }

            const response = await fetch('/api/matching/potential', {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            });

            if (!response.ok) {
                throw new Error('載入配對對象失敗');
            }

            const data = await response.json();
            const newCards = data.data || [];
            
            if (newCards.length === 0) {
                this.showNoMoreCards();
                return;
            }

            this.cards = [...this.cards, ...newCards];
            this.renderCards();
            
            if (this.currentCardIndex === 0) {
                this.showCurrentCard();
            }

        } catch (error) {
            console.error('載入配對對象失敗:', error);
            this.showError('載入配對對象失敗，請稍後再試');
        } finally {
            this.isLoading = false;
            this.showLoading(false);
        }
    }

    renderCards() {
        const container = document.getElementById('cards-container');
        if (!container) return;

        // 只渲染當前卡片和後面的幾張卡片
        const visibleCards = this.cards.slice(this.currentCardIndex, this.currentCardIndex + 3);
        
        container.innerHTML = visibleCards.map((card, index) => {
            const cardIndex = this.currentCardIndex + index;
            return this.createCardHTML(card, cardIndex, index);
        }).join('');

        this.setupCardEvents();
    }

    createCardHTML(user, cardIndex, stackIndex) {
        const age = this.calculateAge(user.birthDate);
        const photos = user.photos || [];
        const primaryPhoto = photos.find(p => p.isPrimary) || photos[0];
        
        return `
            <div class="card" 
                 data-card-index="${cardIndex}"
                 data-stack-index="${stackIndex}"
                 style="z-index: ${10 - stackIndex}; transform: scale(${1 - stackIndex * 0.05}) translateY(${stackIndex * 10}px)">
                
                <div class="card-photos">
                    <div class="photo-indicators">
                        ${photos.map((_, i) => `<div class="indicator ${i === 0 ? 'active' : ''}"></div>`).join('')}
                    </div>
                    <img src="${primaryPhoto?.url || '/static/images/default-avatar.png'}" 
                         alt="${user.name}" class="card-photo">
                    
                    ${photos.length > 1 ? `
                        <div class="photo-nav">
                            <button class="photo-nav-btn prev" onclick="matchingInterface.prevPhoto(${cardIndex})">‹</button>
                            <button class="photo-nav-btn next" onclick="matchingInterface.nextPhoto(${cardIndex})">›</button>
                        </div>
                    ` : ''}
                </div>

                <div class="card-info">
                    <div class="card-header">
                        <h2 class="card-name">${user.name}</h2>
                        <span class="card-age">${age}</span>
                        <div class="card-distance">
                            📍 ${user.location?.city || '未知'} ${this.formatDistance(user.distance)}
                        </div>
                    </div>
                    
                    ${user.bio ? `<p class="card-bio">${user.bio}</p>` : ''}
                    
                    <div class="card-details">
                        ${user.height ? `<div class="detail-item">📏 ${user.height}cm</div>` : ''}
                        ${user.education ? `<div class="detail-item">🎓 ${user.education}</div>` : ''}
                        ${user.occupation ? `<div class="detail-item">💼 ${user.occupation}</div>` : ''}
                    </div>
                    
                    ${user.interests && user.interests.length > 0 ? `
                        <div class="card-interests">
                            ${user.interests.slice(0, 3).map(interest => 
                                `<span class="interest-chip">${interest}</span>`
                            ).join('')}
                            ${user.interests.length > 3 ? `<span class="interest-more">+${user.interests.length - 3}</span>` : ''}
                        </div>
                    ` : ''}
                </div>

                <div class="swipe-indicators">
                    <div class="swipe-indicator swipe-pass">PASS</div>
                    <div class="swipe-indicator swipe-like">LIKE</div>
                    <div class="swipe-indicator swipe-super">SUPER LIKE</div>
                </div>
            </div>
        `;
    }

    setupCardEvents() {
        // 設置拖拽事件到當前卡片
        const currentCard = document.querySelector('.card[data-stack-index="0"]');
        if (currentCard) {
            currentCard.addEventListener('click', (e) => {
                // 阻止在拖拽時觸發點擊
                if (this.isDragging) {
                    e.preventDefault();
                    e.stopPropagation();
                }
            });
        }
    }

    // 觸控事件處理
    handleTouchStart(e) {
        if (e.touches.length !== 1) return;
        
        const touch = e.touches[0];
        this.startDrag(touch.clientX, touch.clientY);
        e.preventDefault();
    }

    handleTouchMove(e) {
        if (!this.isDragging || e.touches.length !== 1) return;
        
        const touch = e.touches[0];
        this.updateDrag(touch.clientX, touch.clientY);
        e.preventDefault();
    }

    handleTouchEnd(e) {
        this.endDrag();
        e.preventDefault();
    }

    // 滑鼠事件處理
    handleMouseDown(e) {
        if (e.button !== 0) return; // 只處理左鍵
        this.startDrag(e.clientX, e.clientY);
        e.preventDefault();
    }

    handleMouseMove(e) {
        if (!this.isDragging) return;
        this.updateDrag(e.clientX, e.clientY);
    }

    handleMouseUp(e) {
        this.endDrag();
    }

    startDrag(x, y) {
        const currentCard = document.querySelector('.card[data-stack-index="0"]');
        if (!currentCard) return;

        this.isDragging = true;
        this.startX = x;
        this.startY = y;
        this.currentX = x;
        this.currentY = y;

        currentCard.classList.add('dragging');
    }

    updateDrag(x, y) {
        if (!this.isDragging) return;

        this.currentX = x;
        this.currentY = y;

        const deltaX = this.currentX - this.startX;
        const deltaY = this.currentY - this.startY;

        const currentCard = document.querySelector('.card[data-stack-index="0"]');
        if (!currentCard) return;

        // 計算旋轉角度
        const rotation = deltaX * this.rotationFactor;
        
        // 應用變換
        currentCard.style.transform = `translate(${deltaX}px, ${deltaY}px) rotate(${rotation}deg)`;
        currentCard.style.opacity = 1 - Math.abs(deltaX) / 300;

        // 顯示滑動指示器
        this.updateSwipeIndicators(deltaX, deltaY);
    }

    endDrag() {
        if (!this.isDragging) return;

        const currentCard = document.querySelector('.card[data-stack-index="0"]');
        if (!currentCard) return;

        const deltaX = this.currentX - this.startX;
        const deltaY = this.currentY - this.startY;

        this.isDragging = false;
        currentCard.classList.remove('dragging');

        // 判斷滑動方向
        if (Math.abs(deltaX) > this.swipeThreshold) {
            if (deltaX > 0) {
                this.executeSwipe('like');
            } else {
                this.executeSwipe('pass');
            }
        } else if (deltaY < -this.swipeThreshold) {
            this.executeSwipe('super-like');
        } else {
            // 回彈到原位
            this.resetCardPosition();
        }

        this.clearSwipeIndicators();
    }

    updateSwipeIndicators(deltaX, deltaY) {
        const passIndicator = document.querySelector('.swipe-pass');
        const likeIndicator = document.querySelector('.swipe-like');
        const superIndicator = document.querySelector('.swipe-super');

        // 重置所有指示器
        [passIndicator, likeIndicator, superIndicator].forEach(indicator => {
            if (indicator) indicator.classList.remove('active');
        });

        // 根據滑動方向顯示對應指示器
        if (Math.abs(deltaX) > this.swipeThreshold) {
            if (deltaX > 0 && likeIndicator) {
                likeIndicator.classList.add('active');
            } else if (deltaX < 0 && passIndicator) {
                passIndicator.classList.add('active');
            }
        } else if (deltaY < -this.swipeThreshold && superIndicator) {
            superIndicator.classList.add('active');
        }
    }

    clearSwipeIndicators() {
        document.querySelectorAll('.swipe-indicator').forEach(indicator => {
            indicator.classList.remove('active');
        });
    }

    resetCardPosition() {
        const currentCard = document.querySelector('.card[data-stack-index="0"]');
        if (!currentCard) return;

        currentCard.style.transform = '';
        currentCard.style.opacity = '1';
    }

    async executeSwipe(action) {
        const currentCard = document.querySelector('.card[data-stack-index="0"]');
        const currentUser = this.cards[this.currentCardIndex];
        
        if (!currentCard || !currentUser) return;

        // 執行動畫
        this.animateCardExit(currentCard, action);

        try {
            // 發送API請求
            const result = await this.sendSwipeAction(currentUser.id, action);
            
            // 檢查是否配對成功
            if (result.isMatch) {
                this.showMatchNotification(currentUser);
            }

        } catch (error) {
            console.error('滑動操作失敗:', error);
            // 不顯示錯誤，但記錄到控制台
        }

        // 移動到下一張卡片
        setTimeout(() => {
            this.nextCard();
        }, 300);
    }

    animateCardExit(card, action) {
        let transform = '';
        
        switch (action) {
            case 'like':
                transform = 'translate(100vw, 0) rotate(30deg)';
                break;
            case 'pass':
                transform = 'translate(-100vw, 0) rotate(-30deg)';
                break;
            case 'super-like':
                transform = 'translate(0, -100vh) rotate(0deg)';
                break;
        }

        card.style.transform = transform;
        card.style.opacity = '0';
        card.style.transition = 'all 0.3s ease-out';
    }

    async sendSwipeAction(userId, action) {
        const token = localStorage.getItem('token');
        
        const response = await fetch('/api/matching/swipe', {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                targetUserId: userId,
                action: action
            })
        });

        if (!response.ok) {
            throw new Error('滑動操作失敗');
        }

        return await response.json();
    }

    showMatchNotification(user) {
        const modal = document.createElement('div');
        modal.className = 'match-modal';
        modal.innerHTML = `
            <div class="match-content">
                <div class="match-animation">
                    <h1 class="match-title">🎉 配對成功！</h1>
                    <div class="match-users">
                        <div class="match-user">
                            <img src="${this.getCurrentUserPhoto()}" alt="您" class="match-photo">
                            <p>您</p>
                        </div>
                        <div class="match-heart">💕</div>
                        <div class="match-user">
                            <img src="${user.photos?.[0]?.url || '/static/images/default-avatar.png'}" alt="${user.name}" class="match-photo">
                            <p>${user.name}</p>
                        </div>
                    </div>
                    <p class="match-message">你們互相喜歡對方！現在可以開始聊天了</p>
                    <div class="match-actions">
                        <button class="btn-secondary" onclick="this.closest('.match-modal').remove()">
                            繼續配對
                        </button>
                        <button class="btn-primary" onclick="window.location.href='/static/html/chat.html'">
                            開始聊天
                        </button>
                    </div>
                </div>
            </div>
        `;

        document.body.appendChild(modal);

        // 自動關閉
        setTimeout(() => {
            if (modal.parentNode) {
                modal.remove();
            }
        }, 10000);
    }

    getCurrentUserPhoto() {
        const user = JSON.parse(localStorage.getItem('user') || '{}');
        return user.photos?.[0]?.url || '/static/images/default-avatar.png';
    }

    nextCard() {
        this.currentCardIndex++;
        
        if (this.currentCardIndex >= this.cards.length - 2) {
            // 當剩餘卡片不多時，載入更多
            this.loadPotentialMatches();
        }
        
        if (this.currentCardIndex >= this.cards.length) {
            this.showNoMoreCards();
            return;
        }

        this.renderCards();
        this.updateCardCounter();
    }

    showNoMoreCards() {
        const container = document.getElementById('cards-container');
        if (container) {
            container.innerHTML = `
                <div class="no-more-cards">
                    <div class="no-cards-icon">💝</div>
                    <h2>沒有更多配對對象了</h2>
                    <p>擴大您的搜尋範圍或稍後再回來看看</p>
                    <div class="no-cards-actions">
                        <button class="btn-secondary" onclick="this.toggleFilters()">調整篩選條件</button>
                        <button class="btn-primary" onclick="location.reload()">重新載入</button>
                    </div>
                </div>
            `;
        }
    }

    updateCardCounter() {
        const counter = document.getElementById('cards-counter');
        if (counter) {
            const remaining = this.cards.length - this.currentCardIndex;
            counter.textContent = remaining > 0 ? `還有 ${remaining} 位` : '沒有更多';
        }
    }

    // 按鈕操作
    swipeLeft() {
        this.executeSwipe('pass');
    }

    swipeRight() {
        this.executeSwipe('like');
    }

    superLike() {
        this.executeSwipe('super-like');
    }

    rewind() {
        if (this.currentCardIndex > 0) {
            this.currentCardIndex--;
            this.renderCards();
            this.updateCardCounter();
        }
    }

    // 照片導航
    nextPhoto(cardIndex) {
        const card = document.querySelector(`[data-card-index="${cardIndex}"]`);
        if (!card) return;

        const photo = card.querySelector('.card-photo');
        const indicators = card.querySelectorAll('.indicator');
        const user = this.cards[cardIndex];
        
        if (!user.photos || user.photos.length <= 1) return;

        let currentPhotoIndex = parseInt(photo.dataset.photoIndex || '0');
        currentPhotoIndex = (currentPhotoIndex + 1) % user.photos.length;
        
        photo.src = user.photos[currentPhotoIndex].url;
        photo.dataset.photoIndex = currentPhotoIndex;

        // 更新指示器
        indicators.forEach((indicator, index) => {
            indicator.classList.toggle('active', index === currentPhotoIndex);
        });
    }

    prevPhoto(cardIndex) {
        const card = document.querySelector(`[data-card-index="${cardIndex}"]`);
        if (!card) return;

        const photo = card.querySelector('.card-photo');
        const indicators = card.querySelectorAll('.indicator');
        const user = this.cards[cardIndex];
        
        if (!user.photos || user.photos.length <= 1) return;

        let currentPhotoIndex = parseInt(photo.dataset.photoIndex || '0');
        currentPhotoIndex = currentPhotoIndex > 0 ? currentPhotoIndex - 1 : user.photos.length - 1;
        
        photo.src = user.photos[currentPhotoIndex].url;
        photo.dataset.photoIndex = currentPhotoIndex;

        // 更新指示器
        indicators.forEach((indicator, index) => {
            indicator.classList.toggle('active', index === currentPhotoIndex);
        });
    }

    // 篩選功能
    toggleFilters() {
        const filtersPanel = document.getElementById('filters-panel');
        if (filtersPanel) {
            filtersPanel.classList.toggle('show');
        }
    }

    applyFilters() {
        // 收集篩選條件
        const filters = {
            ageMin: parseInt(document.getElementById('filter-age-min')?.value || '18'),
            ageMax: parseInt(document.getElementById('filter-age-max')?.value || '65'),
            maxDistance: parseInt(document.getElementById('filter-distance')?.value || '50'),
            interests: Array.from(document.querySelectorAll('#filters-panel .interest-filter:checked'))
                          .map(cb => cb.value)
        };

        // 重置並重新載入
        this.cards = [];
        this.currentCardIndex = 0;
        this.loadPotentialMatches();
        this.toggleFilters();
    }

    // 鍵盤快捷鍵
    handleKeyPress(e) {
        if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') return;

        switch (e.key) {
            case 'ArrowLeft':
                e.preventDefault();
                this.swipeLeft();
                break;
            case 'ArrowRight':
                e.preventDefault();
                this.swipeRight();
                break;
            case 'ArrowUp':
                e.preventDefault();
                this.superLike();
                break;
            case 'Backspace':
                e.preventDefault();
                this.rewind();
                break;
        }
    }

    // 工具函數
    calculateAge(birthDate) {
        if (!birthDate) return '?';
        
        const birth = new Date(birthDate);
        const today = new Date();
        let age = today.getFullYear() - birth.getFullYear();
        const monthDiff = today.getMonth() - birth.getMonth();
        
        if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birth.getDate())) {
            age--;
        }
        
        return age;
    }

    formatDistance(distance) {
        if (!distance) return '';
        return distance < 1000 ? `${Math.round(distance)}m` : `${Math.round(distance / 1000)}km`;
    }

    showCurrentCard() {
        if (this.cards.length > 0) {
            this.renderCards();
            this.updateCardCounter();
        }
    }

    showLoading(show) {
        const loader = document.getElementById('loading-indicator');
        if (loader) {
            loader.style.display = show ? 'flex' : 'none';
        }
    }

    showError(message) {
        const toast = document.createElement('div');
        toast.className = 'toast toast-error';
        toast.innerHTML = `
            <span class="toast-message">${message}</span>
            <button class="toast-close" onclick="this.parentElement.remove()">×</button>
        `;

        document.body.appendChild(toast);

        setTimeout(() => {
            if (toast.parentNode) {
                toast.remove();
            }
        }, 5000);
    }
}

// 全域實例
let matchingInterface;

// 初始化
document.addEventListener('DOMContentLoaded', () => {
    matchingInterface = new MatchingInterface();
});

// 匯出供其他模組使用
if (typeof module !== 'undefined' && module.exports) {
    module.exports = MatchingInterface;
}