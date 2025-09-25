/**
 * 個人檔案管理 JavaScript
 * 包含照片上傳、表單處理、API整合、位置服務
 */

class ProfileManager {
    constructor() {
        this.photos = [];
        this.maxPhotos = 6;
        this.currentUser = null;
        this.isDirty = false; // 追踪是否有未儲存的變更
        
        this.init();
    }

    async init() {
        try {
            await this.loadUserProfile();
            this.bindEvents();
            this.setupValidation();
            this.initializeLocationService();
            this.setupAutoSave();
        } catch (error) {
            console.error('初始化失敗:', error);
            this.showError('載入個人檔案失敗');
        }
    }

    async loadUserProfile() {
        const token = localStorage.getItem('token');
        if (!token) {
            window.location.href = '/static/html/register.html';
            return;
        }

        try {
            const response = await fetch('/api/users/profile', {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            });

            if (!response.ok) {
                throw new Error('載入檔案失敗');
            }

            const data = await response.json();
            this.currentUser = data.data;
            this.populateForm(this.currentUser);
            
        } catch (error) {
            console.error('載入用戶檔案失敗:', error);
            // 如果是認證失敗，重新導向到登入頁面
            if (error.message.includes('401')) {
                localStorage.removeItem('token');
                window.location.href = '/static/html/register.html';
            }
        }
    }

    populateForm(userData) {
        // 基本資訊
        document.getElementById('name').value = userData.name || '';
        document.getElementById('bio').value = userData.bio || '';
        document.getElementById('height').value = userData.height || 170;
        document.getElementById('education').value = userData.education || '';
        document.getElementById('occupation').value = userData.occupation || '';
        
        // 更新滑桿顯示
        this.updateSliderDisplay('height', userData.height || 170);
        
        // 興趣標籤
        if (userData.interests) {
            userData.interests.forEach(interest => {
                const checkbox = document.querySelector(`input[name="interests"][value="${interest}"]`);
                if (checkbox) {
                    checkbox.checked = true;
                    checkbox.parentElement.classList.add('selected');
                }
            });
        }
        
        // 照片
        if (userData.photos) {
            this.photos = userData.photos;
            this.renderPhotos();
        }
        
        // 位置資訊
        if (userData.location) {
            document.getElementById('city').value = userData.location.city || '';
            this.updateLocationDisplay(userData.location);
        }
        
        // 偏好設定
        if (userData.preferences) {
            const prefs = userData.preferences;
            document.getElementById('age-min').value = prefs.ageMin || 18;
            document.getElementById('age-max').value = prefs.ageMax || 65;
            document.getElementById('distance-max').value = prefs.maxDistance || 50;
            
            this.updateSliderDisplay('age-min', prefs.ageMin || 18);
            this.updateSliderDisplay('age-max', prefs.ageMax || 65);
            this.updateSliderDisplay('distance-max', prefs.maxDistance || 50);
        }
    }

    bindEvents() {
        // 照片上傳
        document.getElementById('photo-upload')?.addEventListener('change', (e) => {
            this.handlePhotoUpload(e.target.files);
        });

        // 拖拽上傳
        const uploadArea = document.getElementById('upload-area');
        if (uploadArea) {
            uploadArea.addEventListener('dragover', this.handleDragOver.bind(this));
            uploadArea.addEventListener('drop', this.handleDrop.bind(this));
        }

        // 滑桿事件
        document.querySelectorAll('input[type="range"]').forEach(slider => {
            slider.addEventListener('input', (e) => {
                this.updateSliderDisplay(e.target.id, e.target.value);
                this.markDirty();
            });
        });

        // 興趣標籤選擇
        document.querySelectorAll('.interest-tag').forEach(tag => {
            tag.addEventListener('click', (e) => {
                this.toggleInterest(e.currentTarget);
            });
        });

        // 表單變更追踪
        document.querySelectorAll('input, textarea, select').forEach(element => {
            element.addEventListener('change', () => {
                this.markDirty();
            });
        });

        // 儲存按鈕
        document.getElementById('save-profile')?.addEventListener('click', (e) => {
            e.preventDefault();
            this.saveProfile();
        });

        // 位置服務
        document.getElementById('get-location')?.addEventListener('click', () => {
            this.getCurrentLocation();
        });

        // 頁面離開提醒
        window.addEventListener('beforeunload', (e) => {
            if (this.isDirty) {
                e.preventDefault();
                e.returnValue = '您有未儲存的變更，確定要離開嗎？';
            }
        });
    }

    setupValidation() {
        // 即時驗證
        document.getElementById('bio')?.addEventListener('input', (e) => {
            const remaining = 500 - e.target.value.length;
            const counter = document.getElementById('bio-counter');
            if (counter) {
                counter.textContent = `還可以輸入 ${remaining} 字`;
                counter.style.color = remaining < 50 ? '#e74c3c' : '#666';
            }
        });

        // 名稱驗證
        document.getElementById('name')?.addEventListener('blur', (e) => {
            this.validateName(e.target.value);
        });
    }

    async handlePhotoUpload(files) {
        if (!files || files.length === 0) return;

        const remainingSlots = this.maxPhotos - this.photos.length;
        if (remainingSlots <= 0) {
            this.showError('最多只能上傳6張照片');
            return;
        }

        const filesToProcess = Array.from(files).slice(0, remainingSlots);
        
        for (const file of filesToProcess) {
            if (!this.validatePhotoFile(file)) continue;
            
            try {
                await this.uploadPhoto(file);
            } catch (error) {
                console.error('照片上傳失敗:', error);
                this.showError(`照片上傳失敗: ${error.message}`);
            }
        }
    }

    validatePhotoFile(file) {
        // 檔案類型檢查
        const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp'];
        if (!allowedTypes.includes(file.type)) {
            this.showError('只支援 JPG、PNG、WebP 格式的圖片');
            return false;
        }

        // 檔案大小檢查 (10MB)
        const maxSize = 10 * 1024 * 1024;
        if (file.size > maxSize) {
            this.showError('圖片檔案不能超過10MB');
            return false;
        }

        return true;
    }

    async uploadPhoto(file) {
        const formData = new FormData();
        formData.append('photo', file);

        this.showLoading(true, '上傳照片中...');

        try {
            const token = localStorage.getItem('token');
            const response = await fetch('/api/users/photos', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`
                },
                body: formData
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.message || '上傳失敗');
            }

            const result = await response.json();
            const photoData = {
                id: result.data.id,
                url: result.data.url,
                isPrimary: this.photos.length === 0 // 第一張設為主要照片
            };

            this.photos.push(photoData);
            this.renderPhotos();
            this.showSuccess('照片上傳成功！');

        } finally {
            this.showLoading(false);
        }
    }

    renderPhotos() {
        const container = document.getElementById('photos-grid');
        if (!container) return;

        container.innerHTML = '';

        // 渲染現有照片
        this.photos.forEach((photo, index) => {
            const photoElement = document.createElement('div');
            photoElement.className = `photo-item ${photo.isPrimary ? 'primary' : ''}`;
            photoElement.innerHTML = `
                <img src="${photo.url}" alt="用戶照片 ${index + 1}">
                <div class="photo-overlay">
                    <div class="photo-actions">
                        <button class="photo-action-btn primary-btn" 
                                onclick="profileManager.setPrimaryPhoto(${index})"
                                title="設為主要照片">
                            ⭐
                        </button>
                        <button class="photo-action-btn delete-btn" 
                                onclick="profileManager.deletePhoto(${index})"
                                title="刪除照片">
                            🗑️
                        </button>
                    </div>
                </div>
                ${photo.isPrimary ? '<div class="primary-badge">主要照片</div>' : ''}
            `;
            container.appendChild(photoElement);
        });

        // 添加上傳按鈕（如果還有空位）
        if (this.photos.length < this.maxPhotos) {
            const uploadSlot = document.createElement('div');
            uploadSlot.className = 'photo-upload-slot';
            uploadSlot.innerHTML = `
                <div class="upload-placeholder" onclick="document.getElementById('photo-upload').click()">
                    <div class="upload-icon">📷</div>
                    <div class="upload-text">點擊上傳照片</div>
                    <div class="upload-hint">${this.photos.length}/${this.maxPhotos}</div>
                </div>
            `;
            container.appendChild(uploadSlot);
        }

        // 更新照片計數
        const counter = document.getElementById('photos-counter');
        if (counter) {
            counter.textContent = `${this.photos.length}/${this.maxPhotos}`;
        }
    }

    setPrimaryPhoto(index) {
        if (index < 0 || index >= this.photos.length) return;

        // 移除現有主要照片標記
        this.photos.forEach(photo => photo.isPrimary = false);
        
        // 設定新的主要照片
        this.photos[index].isPrimary = true;
        
        this.renderPhotos();
        this.markDirty();
        this.showSuccess('已設為主要照片');
    }

    async deletePhoto(index) {
        if (index < 0 || index >= this.photos.length) return;

        const photo = this.photos[index];
        if (!confirm('確定要刪除這張照片嗎？')) return;

        try {
            const token = localStorage.getItem('token');
            const response = await fetch(`/api/users/photos/${photo.id}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (!response.ok) {
                throw new Error('刪除照片失敗');
            }

            this.photos.splice(index, 1);
            
            // 如果刪除的是主要照片，設定第一張為主要照片
            if (photo.isPrimary && this.photos.length > 0) {
                this.photos[0].isPrimary = true;
            }
            
            this.renderPhotos();
            this.showSuccess('照片已刪除');

        } catch (error) {
            console.error('刪除照片失敗:', error);
            this.showError('刪除照片失敗');
        }
    }

    handleDragOver(e) {
        e.preventDefault();
        e.currentTarget.classList.add('drag-over');
    }

    handleDrop(e) {
        e.preventDefault();
        e.currentTarget.classList.remove('drag-over');
        this.handlePhotoUpload(e.dataTransfer.files);
    }

    toggleInterest(tagElement) {
        const checkbox = tagElement.querySelector('input[type="checkbox"]');
        if (!checkbox) return;

        const isSelected = checkbox.checked;
        const selectedCount = document.querySelectorAll('.interest-tag input:checked').length;

        // 限制最大選擇數量
        if (!isSelected && selectedCount >= 10) {
            this.showError('最多只能選擇10個興趣');
            return;
        }

        checkbox.checked = !isSelected;
        tagElement.classList.toggle('selected', !isSelected);
        this.markDirty();

        // 更新計數器
        const counter = document.getElementById('interests-counter');
        if (counter) {
            const newCount = document.querySelectorAll('.interest-tag input:checked').length;
            counter.textContent = `已選擇 ${newCount}/10`;
        }
    }

    updateSliderDisplay(sliderId, value) {
        const display = document.getElementById(`${sliderId}-display`);
        if (!display) return;

        let displayText = value;
        
        switch (sliderId) {
            case 'height':
                displayText = `${value} cm`;
                break;
            case 'age-min':
            case 'age-max':
                displayText = `${value} 歲`;
                break;
            case 'distance-max':
                displayText = `${value} km`;
                break;
        }
        
        display.textContent = displayText;
    }

    async getCurrentLocation() {
        if (!navigator.geolocation) {
            this.showError('您的瀏覽器不支援定位功能');
            return;
        }

        this.showLoading(true, '正在取得位置...');

        try {
            const position = await new Promise((resolve, reject) => {
                navigator.geolocation.getCurrentPosition(resolve, reject, {
                    timeout: 10000,
                    enableHighAccuracy: true
                });
            });

            const { latitude, longitude } = position.coords;
            const locationInfo = await this.reverseGeocode(latitude, longitude);
            
            this.updateLocationDisplay(locationInfo);
            document.getElementById('city').value = locationInfo.city || '';
            
            this.markDirty();
            this.showSuccess('位置已更新');

        } catch (error) {
            console.error('獲取位置失敗:', error);
            this.showError('無法獲取您的位置，請手動輸入城市名稱');
        } finally {
            this.showLoading(false);
        }
    }

    async reverseGeocode(lat, lng) {
        // 這裡可以整合真實的地理編碼服務
        // 暫時返回模擬資料
        return {
            city: '台北市',
            district: '信義區',
            latitude: lat,
            longitude: lng
        };
    }

    updateLocationDisplay(location) {
        const display = document.getElementById('location-display');
        if (display && location.city) {
            display.textContent = `📍 ${location.city}${location.district ? ', ' + location.district : ''}`;
            display.style.display = 'block';
        }
    }

    validateName(name) {
        const namePattern = /^[\u4e00-\u9fa5a-zA-Z\s]{2,20}$/;
        const errorElement = document.getElementById('name-error');
        
        if (!name || name.trim().length < 2) {
            this.showFieldError('name', '姓名至少需要2個字符');
            return false;
        }
        
        if (!namePattern.test(name)) {
            this.showFieldError('name', '姓名只能包含中文、英文字母和空格');
            return false;
        }
        
        this.clearFieldError('name');
        return true;
    }

    showFieldError(fieldId, message) {
        const field = document.getElementById(fieldId);
        if (!field) return;

        field.classList.add('error');
        
        let errorElement = document.getElementById(`${fieldId}-error`);
        if (!errorElement) {
            errorElement = document.createElement('div');
            errorElement.id = `${fieldId}-error`;
            errorElement.className = 'field-error';
            field.parentNode.insertBefore(errorElement, field.nextSibling);
        }
        
        errorElement.textContent = message;
        errorElement.style.display = 'block';
    }

    clearFieldError(fieldId) {
        const field = document.getElementById(fieldId);
        if (field) {
            field.classList.remove('error');
        }
        
        const errorElement = document.getElementById(`${fieldId}-error`);
        if (errorElement) {
            errorElement.style.display = 'none';
        }
    }

    async saveProfile() {
        if (!this.validateForm()) {
            this.showError('請修正表單中的錯誤');
            return;
        }

        const profileData = this.collectFormData();
        
        try {
            this.showLoading(true, '儲存中...');
            
            const token = localStorage.getItem('token');
            const response = await fetch('/api/users/profile', {
                method: 'PUT',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(profileData)
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.message || '儲存失敗');
            }

            const result = await response.json();
            this.currentUser = result.data;
            this.markClean();
            
            this.showSuccess('檔案已儲存！');

        } catch (error) {
            console.error('儲存檔案失敗:', error);
            this.showError(error.message || '儲存失敗，請稍後再試');
        } finally {
            this.showLoading(false);
        }
    }

    validateForm() {
        let isValid = true;
        
        // 驗證姓名
        const name = document.getElementById('name').value;
        if (!this.validateName(name)) {
            isValid = false;
        }
        
        // 驗證必填欄位
        const requiredFields = ['name'];
        requiredFields.forEach(fieldId => {
            const field = document.getElementById(fieldId);
            if (!field.value.trim()) {
                this.showFieldError(fieldId, '此欄位為必填');
                isValid = false;
            }
        });
        
        return isValid;
    }

    collectFormData() {
        const data = {
            name: document.getElementById('name').value.trim(),
            bio: document.getElementById('bio').value.trim(),
            height: parseInt(document.getElementById('height').value),
            education: document.getElementById('education').value,
            occupation: document.getElementById('occupation').value,
            location: {
                city: document.getElementById('city').value.trim()
            },
            interests: Array.from(document.querySelectorAll('.interest-tag input:checked'))
                         .map(input => input.value),
            preferences: {
                ageMin: parseInt(document.getElementById('age-min').value),
                ageMax: parseInt(document.getElementById('age-max').value),
                maxDistance: parseInt(document.getElementById('distance-max').value)
            }
        };
        
        return data;
    }

    setupAutoSave() {
        // 每30秒自動儲存
        setInterval(() => {
            if (this.isDirty && this.validateForm()) {
                this.saveProfile();
            }
        }, 30000);
    }

    markDirty() {
        this.isDirty = true;
        this.updateSaveButtonState();
    }

    markClean() {
        this.isDirty = false;
        this.updateSaveButtonState();
    }

    updateSaveButtonState() {
        const saveBtn = document.getElementById('save-profile');
        if (saveBtn) {
            saveBtn.disabled = !this.isDirty;
            saveBtn.textContent = this.isDirty ? '儲存變更' : '已儲存';
        }
    }

    showLoading(show, message = '載入中...') {
        const loadingOverlay = document.getElementById('loading-overlay') || this.createLoadingOverlay();
        const loadingMessage = loadingOverlay.querySelector('.loading-message');
        
        if (show) {
            loadingMessage.textContent = message;
            loadingOverlay.style.display = 'flex';
        } else {
            loadingOverlay.style.display = 'none';
        }
    }

    createLoadingOverlay() {
        const overlay = document.createElement('div');
        overlay.id = 'loading-overlay';
        overlay.className = 'loading-overlay';
        overlay.innerHTML = `
            <div class="loading-content">
                <div class="loading-spinner"></div>
                <div class="loading-message">載入中...</div>
            </div>
        `;
        document.body.appendChild(overlay);
        return overlay;
    }

    showError(message) {
        this.showMessage(message, 'error');
    }

    showSuccess(message) {
        this.showMessage(message, 'success');
    }

    showMessage(message, type) {
        const existingToasts = document.querySelectorAll('.toast');
        existingToasts.forEach(toast => toast.remove());

        const toast = document.createElement('div');
        toast.className = `toast toast-${type}`;
        toast.innerHTML = `
            <div class="toast-content">
                <span class="toast-icon">${type === 'error' ? '❌' : '✅'}</span>
                <span class="toast-message">${message}</span>
            </div>
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
let profileManager;

// 初始化
document.addEventListener('DOMContentLoaded', () => {
    profileManager = new ProfileManager();
});

// 匯出供其他模組使用
if (typeof module !== 'undefined' && module.exports) {
    module.exports = ProfileManager;
}