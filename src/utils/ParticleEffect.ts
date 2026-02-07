export default class ParticleEffect {
  canvas: HTMLCanvasElement;
  ctx: CanvasRenderingContext2D;
  width: number;
  height: number;
  particles: Particle[];
  animationId: number | null = null;
  limit: number;

  constructor(canvas: HTMLCanvasElement, limit: number = 50) {
    this.canvas = canvas;
    const context = this.canvas.getContext('2d');
    if (!context) throw new Error("Could not get canvas context");
    this.ctx = context;
    this.limit = limit;
    this.width = canvas.width = window.innerWidth;
    this.height = canvas.height = window.innerHeight;
    this.particles = [];
    
    this.init();
    window.addEventListener('resize', this.resize.bind(this));
  }

  resize() {
    this.width = this.canvas.width = window.innerWidth;
    this.height = this.canvas.height = window.innerHeight;
  }

  setLimit(limit: number) {
    this.limit = limit;
    if (this.particles.length > limit) {
      this.particles = this.particles.slice(0, limit);
    } else {
      const diff = limit - this.particles.length;
      for (let i = 0; i < diff; i++) {
        this.particles.push(new Particle(this.width, this.height));
      }
    }
  }

  init() {
    this.particles = [];
    for (let i = 0; i < this.limit; i++) {
      this.particles.push(new Particle(this.width, this.height));
    }
  }

  animate() {
    this.ctx.clearRect(0, 0, this.width, this.height);
    this.particles.forEach(p => {
      p.update();
      p.draw(this.ctx);
    });
    this.animationId = requestAnimationFrame(this.animate.bind(this));
  }

  start() {
    if (!this.animationId) {
      this.animate();
    }
  }

  stop() {
    if (this.animationId) {
      cancelAnimationFrame(this.animationId);
      this.animationId = null;
    }
  }
  
  dispose() {
    this.stop();
    window.removeEventListener('resize', this.resize.bind(this));
  }
}

class Particle {
  x: number;
  y: number;
  size: number;
  speedX: number;
  speedY: number;
  opacity: number;
  w: number;
  h: number;

  constructor(w: number, h: number) {
    this.w = w;
    this.h = h;
    this.x = Math.random() * w;
    this.y = Math.random() * h;
    this.size = Math.random() * 2 + 1; // 1-3px
    this.speedX = Math.random() * 0.6 - 0.3; // -0.3 to 0.3
    this.speedY = Math.random() * 0.6 - 0.3;
    this.opacity = Math.random() * 0.5 + 0.2;
  }

  update() {
    this.x += this.speedX;
    this.y += this.speedY;

    if (this.x > this.w) this.x = 0;
    else if (this.x < 0) this.x = this.w;
    
    if (this.y > this.h) this.y = 0;
    else if (this.y < 0) this.y = this.h;
  }

  draw(ctx: CanvasRenderingContext2D) {
    ctx.fillStyle = `rgba(255, 255, 255, ${this.opacity})`;
    ctx.beginPath();
    ctx.arc(this.x, this.y, this.size, 0, Math.PI * 2);
    ctx.fill();
  }
}
