// ECharts 深色主题配置
export const darkTheme = {
  // 全局调色盘
  color: [
    '#fb923c', // 主橙色
    '#4ade80', // 成功绿
    '#60a5fa', // 信息蓝
    '#fbbf24', // 警告黄
    '#f87171', // 危险红
    '#c084fc', // 紫色
    '#2dd4bf', // 青色
    '#f472b6', // 粉色
  ],
  
  // 背景色
  backgroundColor: 'transparent',
  
  // 标题样式
  title: {
    textStyle: {
      color: '#f1f5f9'
    },
    subtextStyle: {
      color: '#94a3b8'
    }
  },
  
  // 图例样式
  legend: {
    textStyle: {
      color: '#cbd5e1'
    },
    pageTextStyle: {
      color: '#cbd5e1'
    }
  },
  
  // 提示框样式
  tooltip: {
    backgroundColor: 'rgba(30, 41, 59, 0.95)',
    borderColor: '#334155',
    borderWidth: 1,
    textStyle: {
      color: '#f1f5f9'
    },
    axisPointer: {
      lineStyle: {
        color: '#475569'
      },
      crossStyle: {
        color: '#475569'
      }
    }
  },
  
  // 坐标轴样式
  axisLine: {
    lineStyle: {
      color: '#334155'
    }
  },
  
  // 坐标轴刻度
  axisTick: {
    lineStyle: {
      color: '#334155'
    }
  },
  
  // 坐标轴标签
  axisLabel: {
    color: '#94a3b8'
  },
  
  // 分割线
  splitLine: {
    lineStyle: {
      color: '#293548',
      type: 'dashed'
    }
  },
  
  // 分割区域
  splitArea: {
    areaStyle: {
      color: ['rgba(255,255,255,0.02)', 'rgba(255,255,255,0)']
    }
  },
  
  // 时间轴
  timeline: {
    lineStyle: {
      color: '#334155'
    },
    itemStyle: {
      color: '#fb923c'
    },
    label: {
      color: '#94a3b8'
    },
    controlStyle: {
      color: '#cbd5e1',
      borderColor: '#475569'
    }
  },
  
  // K线图
  candlestick: {
    itemStyle: {
      color: '#4ade80',
      color0: '#f87171',
      borderColor: '#4ade80',
      borderColor0: '#f87171'
    }
  },
  
  // 仪表盘
  gauge: {
    axisLine: {
      lineStyle: {
        color: [[0.2, '#4ade80'], [0.8, '#fb923c'], [1, '#f87171']]
      }
    },
    axisTick: {
      lineStyle: {
        color: '#475569'
      }
    },
    axisLabel: {
      color: '#94a3b8'
    },
    splitLine: {
      lineStyle: {
        color: '#475569'
      }
    },
    pointer: {
      itemStyle: {
        color: '#fb923c'
      }
    },
    title: {
      color: '#cbd5e1'
    },
    detail: {
      color: '#f1f5f9'
    }
  },
  
  // 数据区域缩放
  dataZoom: {
    backgroundColor: 'rgba(255,255,255,0.02)',
    borderColor: '#334155',
    textStyle: {
      color: '#94a3b8'
    },
    handleStyle: {
      color: '#fb923c'
    },
    moveHandleStyle: {
      color: '#fb923c'
    }
  },
  
  // 标记线
  markLine: {
    lineStyle: {
      color: '#475569'
    },
    label: {
      color: '#cbd5e1'
    }
  },
  
  // 文本样式
  textStyle: {
    color: '#cbd5e1'
  },
  
  // 标签样式
  label: {
    color: '#94a3b8'
  },
  
  // 饼图标签线
  labelLine: {
    lineStyle: {
      color: '#475569'
    }
  }
}

// 应用深色主题的辅助函数
export function applyDarkTheme(echarts) {
  echarts.registerTheme('dark', darkTheme)
}

// 获取主题名称的辅助函数
export function getEChartsTheme() {
  const isDark = document.documentElement.getAttribute('data-theme') === 'dark'
  return isDark ? 'dark' : null
}