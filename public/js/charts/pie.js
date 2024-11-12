import Chart from './base.js';

class AsciiPieChart extends Chart {
    constructor(data, options = {}) {
        const mergedOptions = {
            width: options.width || 19,
            dotChar: '•',
            legendChar: '■',
            colors: ['pie1', 'pie2', 'pie3', 'pie4', 'pie5', 'pie6']
        };

        // Transform object data into array format if needed
        const processedData = Array.isArray(data) 
            ? data 
            : Object.entries(data).map(([label, value]) => ({ label, value }));

        super(processedData, mergedOptions);
        
        this.total = this.values.reduce((sum, val) => sum + val, 0);
        this.segments = this.calculateSegments();
    }

    processValues() {
        return this.data.map(d => {
            const val = typeof d === 'object' ? d.value : d;
            if (typeof val === 'number') return val;
            if (typeof val === 'string') {
                const num = parseFloat(val.replace(/[^0-9.-]/g, ''));
                return isNaN(num) ? 0 : num;
            }
            return 0;
        });
    }

    calculateSegments() {
        let startAngle = 0;
        return this.data.map((item, index) => {
            const percentage = (this.values[index] / this.total) * 100;
            const angleSize = (percentage / 100) * 360;
            const segment = {
                label: item.label,
                value: this.values[index],
                percentage,
                startAngle,
                endAngle: startAngle + angleSize,
                color: this.options.colors[index % this.options.colors.length]
            };
            startAngle += angleSize;
            return segment;
        });
    }

    render() {
        const pieRows = this.generatePieAscii();
        const legendRows = this.generateLegend();
        
        return this.wrapOutput(`
<div class="chart-container">
    <div class="chart-pie">
${pieRows.join('\n')}
    </div>
    <div class="chart-legend">
${legendRows.join('\n')}
    </div>
</div>`);
    }

    generatePieAscii() {
        const rows = [];
        const radius = Math.floor(this.options.width / 2);
        const aspectRatio = 2.1; // Compensate for terminal character height/width ratio
        
        // For each row
        for (let y = -radius; y <= radius; y++) {
            let row = '';
            // For each column
            for (let x = -radius; x <= radius; x++) {
                // Adjust for aspect ratio
                const normalizedX = x;
                const normalizedY = y * aspectRatio;
                
                // Calculate distance from center
                const distance = Math.sqrt(normalizedX * normalizedX + normalizedY * normalizedY);
                
                // If point is within circle radius
                if (distance <= radius) {
                    // Calculate angle in degrees
                    let angle = Math.atan2(normalizedY, normalizedX) * (180 / Math.PI);
                    // Adjust angle to be 0-360
                    angle = angle < 0 ? angle + 360 : angle;
                    
                    // Find which segment this angle belongs to
                    const segment = this.segments.find(seg => 
                        angle >= seg.startAngle && angle < seg.endAngle);
                    
                    row += `<span class="chart-pie-segment" data-palette="${segment.color}">${this.options.dotChar}</span>`;
                } else {
                    row += ' ';
                }
            }
            rows.push(row);
        }
        
        return rows;
    }

    generateLegend() {
        return this.segments.map(segment => {
            const percentage = segment.percentage.toFixed(1).padStart(5);
            const value = segment.value.toString().padStart(4);
            return `<span class="chart-legend-item" data-palette="${segment.color}">${this.options.legendChar}</span> ${percentage}% [${value}] ${segment.label}`;
        });
    }

    wrapOutput(content) {
        return `<pre class="chart">${content}</pre>`;
    }
}

export default AsciiPieChart; 