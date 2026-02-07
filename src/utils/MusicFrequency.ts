// 用于在 Canvas 上绘制音乐频谱
class MusicFrequency {
  canvas: HTMLCanvasElement;
  audio: HTMLAudioElement;
  width: number;
  height: number;
  color: string;
  lineWidth: number;
  vHight: number;
  ctx: CanvasRenderingContext2D;
  output: Uint8Array;
  context!: AudioContext;
  source!: MediaElementAudioSourceNode;
  analyser!: AnalyserNode;
  grd!: CanvasGradient;
  animationId: number | null = null;
  private isDispose: boolean = false;
  scale: number = 0.9;

  /**
   * 创建一个 MusicFrequency 实例。
   * @param {HTMLCanvasElement} canvas - 要绘制频谱的画布元素。
   * @param {HTMLAudioElement} audio - 要播放并可视化的音频元素。
   * @param {string} [color] - 主题颜色 hex.
   * @param {number} [width] - 画布的宽度，默认为页面宽度的 1600 分之一或当前窗口的宽度（取较小值）。
   * @param {number} [height] - 画布的高度，默认为 200。
   * @param {number} [lineWidth] - 线条宽度，默认为画布宽度的 360 分之一的 1.6 倍。
   * @param {number} [vHight] - 纵向缩放比例 - 已废弃，内部计算以适应高度。
   * @param {number} [scale] - 频谱跳动幅度 (0-200) -> 映射到 (0-2.0)。
   */
  constructor(
    canvas: HTMLCanvasElement,
    audio: HTMLAudioElement,
    color: string = "#ffffff",
    width: number | null = null,
    height: number | null = null,
    lineWidth: number | null = null,
    vHight: number | null = null,
    scale: number = 90
  ) {
    this.canvas = canvas;
    // 设置画布宽高
    this.width = width || (document.body.clientWidth >= 1600 ? 1600 : document.body.clientWidth);
    this.canvas.width = this.width;
    this.height = height || 200;
    this.canvas.height = this.height;
    
    // 存储音频和其他参数
    this.audio = audio;
    this.color = color;
    // 增加高度差，这里稍微减小默认除数，或者使用指数
    this.vHight = vHight || 2;
    this.scale = scale / 100;

    // 计算线条宽度
    this.lineWidth = lineWidth || this.canvas.width / 360 / 1.6;
    
    // 获取画布上下文
    const context = this.canvas.getContext("2d");
    if (!context) throw new Error("Could not get canvas context");
    this.ctx = context;
    
    // 创建渐变
    this.createGradient();
    
    // 创建输出缓冲区
    this.output = new Uint8Array(360);
    
    // 劫持 audio.play 和 pause，确保 AudioContext 状态正确
    // 注意：这样修改可能会有副作用，通常建议外部管理 AudioContext，但为了保持兼容性保留逻辑
    const originalPlay = this.audio.play.bind(this.audio);
    this.audio.play = async () => {
      if (this.context && this.context.state === 'suspended') {
        await this.context.resume();
      }
      return originalPlay();
    };
    
    /* 
       不建议覆盖 pause 导致 suspend，因为暂停时可能还需要其他音频处理
       或者仅仅是为了省资源。保留原逻辑。
    */
   
    // 创建音频上下文和分析器
    this.createContext();
  }

  // 设置颜色并更新
  setColor(color: string) {
    this.color = color;
    this.createGradient();
  }

   // 设置幅度
   setScale(scale: number) {
      this.scale = scale / 100;
   }


  // 创建渐变
  createGradient() {
    this.grd = this.ctx.createLinearGradient(0, this.canvas.height, 0, 0);
    
    // 使用 hexToRgba 转换颜色以实现透明度渐变
    // 底部稍微透明，顶部不透明
    const hex = this.color; 
    
    this.grd.addColorStop(0, `${hex}00`); // 底部透明
    this.grd.addColorStop(0.5, `${hex}80`); // 中间半透明
    this.grd.addColorStop(1, `${hex}ff`);   // 顶部完全不透明
    
    // 或者按照用户之前的逻辑：
    // 如果想要纯色渐变，可以都设为 this.color，或者根据用户需求 "颜色跟随主题色"
    // 之前的代码是白色到白色。现在可以是 主题色到主题色(带透明度变化会好看些)
  }

  // 创建音频上下文和分析器
  createContext() {
    // 创建音频上下文
    const AudioContext = window.AudioContext || (window as any).webkitAudioContext;
    if (!this.context) {
        this.context = new AudioContext();
    }
    
    // 防止重复连接
    if (!this.source) {
      // 创建媒体源和分析器
      try {
        this.source = this.context.createMediaElementSource(this.audio);
      } catch (e) {
        // 如果已经连接过 source，再次创建会报错，这里忽略或者重用
        console.warn("MediaElementAudioSourceNode create failed:", e);
        return; 
      }
      this.analyser = this.context.createAnalyser();
      this.analyser.fftSize = 2048; // 增加 FFT 大小以获得更准确的数据
      // 连接媒体源和分析器
      this.source.connect(this.analyser);
      this.analyser.connect(this.context.destination);
    }
  }

  drawSpectrum() {
    if (this.isDispose) return;

    // 获取频域数据
    this.analyser.getByteFrequencyData(this.output as any);
    // 清除画布
    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
    
    // 绘制每个频率的柱形
    // 360 个点可能对应 output 的前 360 个数据，还是取样？
    // 原代码直接使用 output (size 360)，但 analyser.fftSize 默认为 2048， frequencyBinCount 是 1024
    // 之前的 output = new Uint8Array(360) 意味着 getByteFrequencyData 只会填充前 360 个 bin
    // 这通常对应低频部分，这是合理的。

    for (let i = 0; i < 360; i++) {
        // 优化高低差：使用幂函数来拉大差异
        // 原逻辑: value = output[i] / vHight
        
        let value = this.output[i] || 0;
        
        // 放大差异算法：
        // 1. 归一化 0-1
        let ratio = value / 255;
        // 2. 指数放大 (立方或者更高) 让小值更小，大值更大
        // 3. 重新放大回像素高度
        
        // 比如使用平方
        const amplified = Math.pow(ratio, 2.5); // 指数越大，差异越明显
        
        // 计算高度，最大高度为 canvas.height
        // 这里的 scaleFactor 可以调整整体高度
        // 如果 vHight 是缩放因子（越小越高），我们动态调整
        
        // 原来是 / 5 (当 height=50 时) => max 255 / 5 = 51 ≈ 50
        // 我们希望保持最大值接近 canvas height
        
        let barHeight = amplified * this.canvas.height * (255 / (this.vHight * 40)); 
        // 这里的系数可能需要微调，为了安全起见，我们直接映射到 height
        
        // 使用简单的线性映射配合指数增强
        // 假设 vHight 仍然控制整体缩放 (兼容旧参数逻辑)
        // 旧: 255 / 5 = 50.
        // 新: 1.0 * H.
        
        barHeight = amplified * (this.canvas.height * this.scale); //留一点顶部空间

        let x = i * (this.width / 360); // 确保填满宽度
        
        if (x <= this.canvas.width) {
            this.ctx.fillStyle = this.grd; // 使用 fillStyle 填充矩形可能更好看？原代码是 stroke 只有线条
            // 如果原代码是线条：
            this.ctx.strokeStyle = this.grd;
            this.ctx.beginPath();
            this.ctx.lineWidth = this.lineWidth;
            this.ctx.lineCap = "round"; // 圆头更好看
            this.ctx.moveTo(x, this.canvas.height);
            this.ctx.lineTo(x, this.canvas.height - barHeight);
            this.ctx.stroke();
        }
    }

    //请求下一帧
    this.animationId = requestAnimationFrame(() => {
      this.drawSpectrum();
    });
  }
  
  dispose() {
      this.isDispose = true;
      if (this.animationId) {
          cancelAnimationFrame(this.animationId);
      }
  }
}

export default MusicFrequency;
