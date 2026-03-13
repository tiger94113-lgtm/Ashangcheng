// ============================================
// Asset Mall - Go 后端 API 客户端
// 将此代码添加到前端 app.js 中使用
// ============================================

const API_CONFIG = {
  // 修改为你的后端地址
  baseURL: "http://localhost:8080/api",
  // 或者使用环境变量
  // baseURL: process.env.REACT_APP_API_URL || "http://localhost:8080/api",
};

// 封装 fetch 请求
async function apiRequest(endpoint, options = {}) {
  const url = `${API_CONFIG.baseURL}${endpoint}`;
  
  const defaultOptions = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  const response = await fetch(url, { ...defaultOptions, ...options });
  
  if (!response.ok) {
    throw new Error(`API Error: ${response.status}`);
  }

  return response.json();
}

// ============================================
// 订单 API
// ============================================

const OrderAPI = {
  // 获取订单列表
  async getOrders(params = {}) {
    const query = new URLSearchParams(params).toString();
    const endpoint = query ? `/orders?${query}` : "/orders";
    return apiRequest(endpoint);
  },

  // 获取单个订单
  async getOrder(id) {
    return apiRequest(`/orders/${id}`);
  },

  // 创建订单
  async createOrder(order) {
    return apiRequest("/orders", {
      method: "POST",
      body: JSON.stringify(order),
    });
  },

  // 更新订单
  async updateOrder(id, updates) {
    return apiRequest(`/orders/${id}`, {
      method: "PUT",
      body: JSON.stringify(updates),
    });
  },

  // 删除订单
  async deleteOrder(id) {
    return apiRequest(`/orders/${id}`, {
      method: "DELETE",
    });
  },

  // 获取统计
  async getStats() {
    return apiRequest("/stats");
  },

  // 健康检查
  async healthCheck() {
    return apiRequest("/health");
  },
};

// ============================================
// 与现有前端集成
// ============================================

// 替换原有的本地存储函数
async function saveOrder(order) {
  try {
    // 先保存到后端
    const result = await OrderAPI.createOrder(order);
    
    // 同时保存到本地作为备份
    const orders = loadStoredOrders();
    orders.unshift(order);
    saveStoredOrders(orders);
    
    return result;
  } catch (error) {
    console.error("保存到后端失败，仅保存到本地:", error);
    // 后端失败时，只保存到本地
    const orders = loadStoredOrders();
    orders.unshift(order);
    saveStoredOrders(orders);
    throw error;
  }
}

// 从后端加载订单
async function loadOrdersFromBackend(wallet) {
  try {
    const result = await OrderAPI.getOrders({ wallet });
    return result.data || [];
  } catch (error) {
    console.error("从后端加载失败:", error);
    // 后端失败时，从本地加载
    return loadStoredOrders().filter(o => o.wallet === wallet);
  }
}

// 同步本地订单到后端
async function syncOrdersToBackend() {
  const localOrders = loadStoredOrders();
  const results = { success: 0, failed: 0 };

  for (const order of localOrders) {
    try {
      await OrderAPI.createOrder(order);
      results.success++;
    } catch (error) {
      results.failed++;
    }
  }

  return results;
}

// ============================================
// 使用示例
// ============================================

/*
// 1. 创建订单时
const order = {
  id: generateUUID(),
  orderNo: createOrderNumber(),
  status: "pending",
  wallet: state.account,
  name: customer.name,
  phone: customer.phone,
  address: customer.address,
  note: customer.note,
  referrer: referrer,
  usdtAmount: usdtAmount.toString(),
  rewardPoolAmount: rewardPoolAmount.toString(),
  assetAmount: assetAmount.toString(),
  minAssetAmount: minAssetAmount.toString(),
  txHash: tx.hash,
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString(),
};

try {
  await saveOrder(order);
  setStatus("订单已保存到服务器", "ok");
} catch (error) {
  setStatus("订单已保存到本地", "warn");
}

// 2. 加载订单时
async function renderOrders() {
  let orders;
  if (state.account) {
    orders = await loadOrdersFromBackend(state.account);
  } else {
    orders = loadStoredOrders();
  }
  // ... 渲染订单
}

// 3. 管理员查看所有订单
async function loadAllOrders() {
  const result = await OrderAPI.getOrders();
  return result.data;
}

// 4. 搜索订单
async function searchOrders(keyword) {
  const result = await OrderAPI.getOrders({ search: keyword });
  return result.data;
}

// 5. 更新订单状态
async function updateOrderStatus(orderId, status) {
  await OrderAPI.updateOrder(orderId, { status });
}

// 6. 获取统计
async function updateStats() {
  const stats = await OrderAPI.getStats();
  console.log("总订单:", stats.totalOrders);
  console.log("已确认:", stats.confirmedOrders);
  console.log("总USDT:", stats.totalUsdt);
}
*/

// 导出 API
window.OrderAPI = OrderAPI;
window.saveOrder = saveOrder;
window.loadOrdersFromBackend = loadOrdersFromBackend;
window.syncOrdersToBackend = syncOrdersToBackend;
