<template>
  <div class="playground-container">
    <div class="toolbar">
      <div class="version-selector">
        <select v-model="selectedVersion">
          <option value="go1.24">Go 1.24</option>
          <option value="go1.23">Go 1.23</option>
          <option value="go1.22">Go 1.22</option>
        </select>
      </div>
      <div class="action-buttons">
        <button @click="runCode" class="btn btn-run" :disabled="isRunning || isFormatting">
          <i class="fas fa-play"></i> {{ isRunning ? '运行中...' : '运行' }}
        </button>
        <button @click="formatCode" class="btn btn-format" :disabled="isRunning || isFormatting">
          <i class="fas fa-align-left"></i> {{ isFormatting ? '格式化中...' : '格式化' }}
        </button>
        <button @click="showShareDialog" class="btn btn-share">
          <i class="fas fa-share-alt"></i> Share
        </button>
      </div>
      <div class="examples-dropdown">
        <select v-model="selectedExample" @change="loadExample">
          <option value="hello">Hello, World!</option>
          <option value="conway">Conway's Game of Life</option>
          <option value="fibonacci">Fibonacci</option>
          <option value="concurrent-pi">Concurrent Pi</option>
          <option value="http-server">HTTP Server</option>
          <option value="goversion">Go Version</option>
        </select>
      </div>
    </div>

    <div class="editor-container">
      <textarea 
        class="code-editor" 
        v-model="code"
        placeholder="// Write your Go code here"
        @keydown.tab.prevent="handleTab"
        spellcheck="false"
        autocorrect="off"
        autocapitalize="off"
        autocomplete="off"
      ></textarea>
    </div>

    <div class="output-container" v-if="output">
      <h3>Output:</h3>
      <pre>{{ output }}</pre>
    </div>

    <!-- 分享对话框 -->
    <div v-if="showShare" class="modal-overlay" @click.self="showShare = false">
      <div class="modal">
        <Share
          :code="code"
          :version="selectedVersion"
          @close="showShare = false"
        />
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import Share from './Share.vue'

export default {
  name: 'Playground',
  components: {
    Share
  },
  data() {
    return {
      code: '// Hello World example\npackage main\n\nimport "fmt"\n\nfunc main() {\n\tfmt.Println("Hello, 世界")\n}',
      output: '',
      selectedVersion: 'go1.24',
      selectedExample: 'hello',
      isRunning: false,
      isFormatting: false,
      examples: {
        'hello': '// Hello World example\npackage main\n\nimport "fmt"\n\nfunc main() {\n\tfmt.Println("Hello, 世界")\n}',
        'fibonacci': '// Fibonacci example\npackage main\n\nimport "fmt"\n\nfunc fibonacci(n int) int {\n\tif n <= 1 {\n\t\treturn n\n\t}\n\treturn fibonacci(n-1) + fibonacci(n-2)\n}\n\nfunc main() {\n\tfor i := 0; i < 10; i++ {\n\t\tfmt.Println(fibonacci(i))\n\t}\n}',
        'conway': '// Conway\'s Game of Life\npackage main\n\nimport (\n\t"bytes"\n\t"fmt"\n)\n\nfunc main() {\n\t// Initialize the world\n\tworld := NewWorld(10, 10)\n\tworld.Set(3, 3, true)\n\tworld.Set(3, 4, true)\n\tworld.Set(3, 5, true)\n\n\t// Run simulation for 5 generations\n\tfor i := 0; i < 5; i++ {\n\t\tfmt.Printf("Generation %d:\\n", i)\n\t\tfmt.Println(world)\n\t\tworld = world.Next()\n\t}\n}\n\ntype World struct {\n\tcells [][]bool\n\twidth, height int\n}\n\nfunc NewWorld(width, height int) *World {\n\tcells := make([][]bool, height)\n\tfor i := range cells {\n\t\tcells[i] = make([]bool, width)\n\t}\n\treturn &World{cells: cells, width: width, height: height}\n}\n\nfunc (w *World) Set(x, y int, alive bool) {\n\tw.cells[y][x] = alive\n}\n\nfunc (w *World) Alive(x, y int) bool {\n\tx = (x + w.width) % w.width\n\ty = (y + w.height) % w.height\n\treturn w.cells[y][x]\n}\n\nfunc (w *World) Next() *World {\n\tnext := NewWorld(w.width, w.height)\n\tfor y := 0; y < w.height; y++ {\n\t\tfor x := 0; x < w.width; x++ {\n\t\t\talive := w.Alive(x, y)\n\t\t\tneighbors := w.AliveNeighbors(x, y)\n\t\t\tif alive && (neighbors < 2 || neighbors > 3) {\n\t\t\t\tnext.Set(x, y, false)\n\t\t\t} else if !alive && neighbors == 3 {\n\t\t\t\tnext.Set(x, y, true)\n\t\t\t} else {\n\t\t\t\tnext.Set(x, y, alive)\n\t\t\t}\n\t\t}\n\t}\n\treturn next\n}\n\nfunc (w *World) AliveNeighbors(x, y int) int {\n\tcount := 0\n\tfor dy := -1; dy <= 1; dy++ {\n\t\tfor dx := -1; dx <= 1; dx++ {\n\t\t\tif dx == 0 && dy == 0 {\n\t\t\t\tcontinue\n\t\t\t}\n\t\t\tif w.Alive(x+dx, y+dy) {\n\t\t\t\tcount++\n\t\t\t}\n\t\t}\n\t}\n\treturn count\n}\n\nfunc (w *World) String() string {\n\tvar buf bytes.Buffer\n\tfor y := 0; y < w.height; y++ {\n\t\tfor x := 0; x < w.width; x++ {\n\t\t\tif w.Alive(x, y) {\n\t\t\t\tbuf.WriteString("O ")\n\t\t\t} else {\n\t\t\t\tbuf.WriteString(". ")\n\t\t\t}\n\t\t}\n\t\tbuf.WriteString("\\n")\n\t}\n\treturn buf.String()\n}',
        'concurrent-pi': '// Concurrent Pi calculation\npackage main\n\nimport (\n\t"fmt"\n\t"math"\n)\n\nfunc main() {\n\tch := make(chan float64)\n\tgo calculatePi(ch, 100000)\n\tpi := <-ch\n\tfmt.Printf("Calculated Pi: %.10f\\n", pi)\n\tfmt.Printf("Actual Pi:     %.10f\\n", math.Pi)\n}\n\nfunc calculatePi(ch chan float64, iterations int) {\n\tvar sum float64 = 0\n\tfor i := 0; i < iterations; i++ {\n\t\tx := (float64(i) + 0.5) / float64(iterations)\n\t\tsum += 4 / (1 + x*x)\n\t}\n\tch <- sum / float64(iterations)\n}',
        'http-server': '// HTTP Server example\npackage main\n\nimport (\n\t"fmt"\n\t"net/http"\n)\n\nfunc main() {\n\thttp.HandleFunc("/", handler)\n\tfmt.Println("Starting server at port 8080...")\n\t// Note: This won\'t actually run in the playground\n\t// as network operations are restricted\n\t// http.ListenAndServe(":8080", nil)\n\tfmt.Println("Server example - network operations disabled")\n}\n\nfunc handler(w http.ResponseWriter, r *http.Request) {\n\tfmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])\n}',
        'goversion': '// Go version example\npackage main\n\nimport (\n\t"fmt"\n\t"runtime"\n)\n\nfunc main() {\n\tfmt.Println("Go version:", runtime.Version())\n}'
      },
      showShare: false
    }
  },
  methods: {
    runCode() {
      if (this.isRunning || this.isFormatting) return;
      this.isRunning = true;
      // Call the backend API
      const versionPath = this.getVersionPath();
      axios.post(`/api/${versionPath}/run`, {
        code: this.code,
        version: this.selectedVersion
      })
      .then(response => {
        // Display output, and if there's an error, display that too
        let result = response.data.output || '';
        if (response.data.error) {
          result += '\n\n-- Error --\n' + response.data.error;
        }
        this.output = result;
      })
      .catch(error => {
        this.output = 'Error: ' + error.response?.data?.error || 'Failed to run code';
      })
      .finally(() => {
        this.isRunning = false;
      });
    },
    formatCode() {
      if (this.isRunning || this.isFormatting) return;
      this.isFormatting = true;
      // Call the backend API for formatting
      const versionPath = this.getVersionPath();
      axios.post(`/api/${versionPath}/format`, {
        code: this.code
      })
      .then(response => {
        this.code = response.data.formattedCode;
      })
      .catch(error => {
        this.output = 'Error formatting: ' + error.response?.data?.error || 'Failed to format code';
      })
      .finally(() => {
        this.isFormatting = false;
      });
    },
    getVersionPath() {
      // Map the selected version to the API path
      switch(this.selectedVersion) {
        case 'go1.24':
          return 'go1.24';
        case 'go1.23':
          return 'go1.23';
        case 'go1.22':
          return 'go1.22';
        default:
          return 'go1.24'; // Default to Go 1.24
      }
    },
    showShareDialog() {
      this.showShare = true
    },
    loadExample() {
      if (this.selectedExample && this.examples[this.selectedExample]) {
        this.code = this.examples[this.selectedExample];
      }
    },
    handleTab(e) {
      // Insert a tab character when tab is pressed
      const start = e.target.selectionStart;
      const end = e.target.selectionEnd;
      
      this.code = this.code.substring(0, start) + '\t' + this.code.substring(end);
      
      // Move cursor position
      this.$nextTick(() => {
        e.target.selectionStart = e.target.selectionEnd = start + 1;
      });
    }
  }
}
</script>

<style>
.playground-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 100%;
  border: 1px solid #ddd;
  border-radius: 4px;
  overflow: hidden;
  flex: 1;
}

.toolbar {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  padding: 0.5rem;
  background-color: #f5f5f5;
  border-bottom: 1px solid #ddd;
  gap: 8px;
}

.version-selector select,
.examples-dropdown select {
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  width: 100%;
  min-width: 120px;
}

.action-buttons {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
  white-space: nowrap;
  display: flex;
  align-items: center;
  gap: 5px;
}

.btn-run {
  background-color: #00ADD8;
  color: white;
}

.btn-format {
  background-color: #E0EBF5;
  color: #375EAB;
}

.btn-share {
  background-color: #375EAB;
  color: white;
}

.editor-container {
  flex: 1;
  min-height: 200px;
  display: flex;
  flex-direction: column;
}

.code-editor {
  width: 100%;
  height: 100%;
  min-height: 200px;
  padding: 1rem;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
  border: none;
  resize: none;
  tab-size: 4;
  flex: 1;
}

.output-container {
  padding: 1rem;
  background-color: #f8f8f8;
  border-top: 1px solid #ddd;
  overflow-x: auto;
  max-height: none;
  width: 100%;
}

.output-container h3 {
  margin-top: 0;
  margin-bottom: 0.5rem;
  font-size: 1rem;
  color: #333;
}

.output-container pre {
  margin: 0;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  white-space: pre-wrap;
  word-break: break-word;
  overflow-x: auto;
  padding-bottom: 0.5rem;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.modal {
  background: white;
  border-radius: 4px;
  width: 100%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
}

/* 媒体查询以适应移动设备 */
@media (max-width: 768px) {
  .playground-container {
    border: none;
    border-radius: 0;
  }
  
  .toolbar {
    flex-direction: column;
    align-items: stretch;
  }
  
  .version-selector,
  .action-buttons,
  .examples-dropdown {
    width: 100%;
    margin-bottom: 8px;
  }
  
  .action-buttons {
    justify-content: space-between;
  }
  
  .btn {
    padding: 0.5rem;
    flex: 1;
    justify-content: center;
  }
  
  .code-editor {
    font-size: 16px; /* 在移动设备上增大字体 */
    padding: 0.75rem;
  }
  
  .output-container {
    padding: 0.75rem;
  }
  
  .output-container h3 {
    font-size: 0.9rem;
  }
  
  .output-container pre {
    font-size: 16px;
    line-height: 1.4;
  }
}

@media (max-width: 480px) {
  .toolbar {
    padding: 0.25rem;
  }
  
  .btn {
    padding: 0.5rem 0.25rem;
    font-size: 13px;
  }
  
  .modal {
    max-width: 100%;
    margin: 0 10px;
  }
  
  .output-container {
    padding: 0.5rem;
    min-height: 120px; /* 确保在小屏幕上输出区域有足够高度 */
  }
  
  .output-container pre {
    padding-bottom: 1rem; /* 增加底部留白 */
  }
}
</style> 