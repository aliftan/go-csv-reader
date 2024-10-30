// static/js/app.js
const { createApp } = Vue

const app = createApp({
    data() {
        return {
            headers: [],
            selectedColumns: [],
            results: [],
            filter: {
                column: '',
                operator: 'equals',
                value: ''
            },
            selectedFile: null,
            originalFileName: ''
        }
    },
    methods: {
        handleFileUpload(event) {
            this.selectedFile = event.target.files[0]
            // Store the original filename without extension
            if (this.selectedFile) {
                this.originalFileName = this.selectedFile.name.replace(/\.[^/.]+$/, '')
            }
        },
        async uploadFile() {
            if (!this.selectedFile) {
                alert('Please select a file')
                return
            }

            const formData = new FormData()
            formData.append('csvFile', this.selectedFile)

            try {
                const response = await fetch('/api/upload', {
                    method: 'POST',
                    body: formData
                })

                if (!response.ok) {
                    throw new Error('Upload failed')
                }

                await this.fetchHeaders()
                await this.queryData()
            } catch (error) {
                alert('Error uploading file: ' + error.message)
            }
        },
        resetUpload() {
            // Reset file input
            const fileInput = document.querySelector('input[type="file"]')
            if (fileInput) {
                fileInput.value = ''
            }

            // Reset all data
            this.selectedFile = null
            this.headers = []
            this.selectedColumns = []
            this.results = []
            this.resetFilter()
        },
        resetFilter() {
            this.filter = {
                column: '',
                operator: 'equals',
                value: ''
            }
            // If we have data, requery without filters
            if (this.headers.length) {
                this.queryData()
            }
        },
        async fetchHeaders() {
            try {
                const response = await fetch('/api/headers')
                this.headers = await response.json()
                this.selectedColumns = [...this.headers]
            } catch (error) {
                alert('Error fetching headers: ' + error.message)
            }
        },
        exportToJSON() {
            try {
                const jsonData = JSON.stringify(this.results, null, 2)
                const blob = new Blob([jsonData], { type: 'application/json' })
                const url = window.URL.createObjectURL(blob)

                // Create filename with timestamp and original filename
                const timestamp = new Date().toISOString().split('T')[0]
                const filename = this.originalFileName
                    ? `${timestamp}_export_${this.originalFileName}.json`
                    : `${timestamp}_export.json`

                this.downloadFile(url, filename)
            } catch (error) {
                alert('Error exporting JSON: ' + error.message)
            }
        },
        exportToCSV() {
            try {
                let csv = this.selectedColumns.join(',') + '\n'

                this.results.forEach(record => {
                    const row = this.selectedColumns.map(column => {
                        const value = record[column] || ''
                        return `"${value.replace(/"/g, '""')}"`
                    })
                    csv += row.join(',') + '\n'
                })

                const blob = new Blob([csv], { type: 'text/csv' })
                const url = window.URL.createObjectURL(blob)

                // Create filename with timestamp and original filename
                const timestamp = new Date().toISOString().split('T')[0]
                const filename = this.originalFileName
                    ? `${timestamp}_export_${this.originalFileName}.csv`
                    : `${timestamp}_export.csv`

                this.downloadFile(url, filename)
            } catch (error) {
                alert('Error exporting CSV: ' + error.message)
            }
        },
        downloadFile(url, filename) {
            const link = document.createElement('a')
            link.href = url
            link.download = filename
            document.body.appendChild(link)
            link.click()
            document.body.removeChild(link)
            window.URL.revokeObjectURL(url)
        },
        async queryData() {
            try {
                const response = await fetch('/api/query', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        columns: this.selectedColumns,
                        filter: this.filter
                    })
                })

                if (!response.ok) {
                    throw new Error('Query failed')
                }

                this.results = await response.json()
            } catch (error) {
                alert('Error querying data: ' + error.message)
            }
        }
    }
})

app.mount('#app')