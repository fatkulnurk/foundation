# Storage Package Examples

This directory contains working examples demonstrating how to use the storage package.

## Prerequisites

1. **Go 1.25 or higher**
2. **For Local Storage:** No additional requirements
3. **For S3 Storage:** AWS account or MinIO server

## Running the Examples

### 1. Local Storage Example

Demonstrates basic local filesystem storage operations.

```bash
cd pkg/storage/example
go run main.go local
```

**What it does:**
- Uploads text and JSON files
- Checks if files exist
- Gets file information
- Reads file content

**Output location:** `./storage/` directory

### 2. S3 Storage Example

Demonstrates AWS S3 storage operations.

**Setup environment variables:**
```bash
export STORAGE_S3_REGION="us-east-1"
export STORAGE_S3_BUCKET="your-bucket-name"
export STORAGE_S3_ACCESS_KEY="your-access-key"
export STORAGE_S3_SECRET_KEY="your-secret-key"
```

**Run:**
```bash
go run main.go s3
```

**What it does:**
- Uploads public and private files to S3
- Generates temporary URLs for private files
- Shows file information

### 3. List Files Example

Demonstrates listing files and directories.

```bash
go run main.go list
```

**What it does:**
- Creates test files in multiple directories
- Lists files in each directory
- Lists subdirectories
- Shows file sizes and MIME types

### 4. File Operations Example

Demonstrates copy, move, and delete operations.

```bash
go run main.go operations
```

**What it does:**
- Creates an original file
- Copies file to backup location
- Moves file to archive location
- Deletes backup file
- Verifies each operation

## Example Output

### Local Storage Example

```
=== Local Storage Example ===

1. Uploading text file...
✓ Uploaded: documents/hello.txt (22 B)

2. Uploading JSON file...
✓ Uploaded: data/user.json (45 B)

3. Checking if file exists...
✓ File exists!

4. Getting file information...
✓ File Info:
  - Name: hello.txt
  - Size: 22 B
  - Type: text/plain
  - URL: http://localhost:8080/storage/documents/hello.txt
  - Modified: 2024-11-26 19:45:30

5. Reading file content...
✓ Content: Hello, Local Storage!

✅ Local storage example completed!
Files are stored in: ./storage/
```

### List Files Example

```
=== List Files Example ===

Creating test files...
✓ Test files created

1. Files in photos/:
  - image1.jpg (12 B)
  - image2.jpg (12 B)

2. Files in documents/:
  - report.pdf (8 B) - application/octet-stream
  - invoice.pdf (12 B) - application/octet-stream

3. Subdirectories in photos/:
  - photos/vacation/

✅ List example completed!
```

## Using MinIO (S3-Compatible)

If you want to test S3 functionality locally, use MinIO:

### 1. Start MinIO with Docker

```bash
docker run -p 9000:9000 -p 9001:9001 \
  -e "MINIO_ROOT_USER=minioadmin" \
  -e "MINIO_ROOT_PASSWORD=minioadmin" \
  minio/minio server /data --console-address ":9001"
```

### 2. Create a Bucket

Open http://localhost:9001 and create a bucket named "test-bucket"

### 3. Set Environment Variables

```bash
export STORAGE_S3_REGION="us-east-1"
export STORAGE_S3_BUCKET="test-bucket"
export STORAGE_S3_ACCESS_KEY="minioadmin"
export STORAGE_S3_SECRET_KEY="minioadmin"
export STORAGE_S3_URL="http://localhost:9000"
export STORAGE_S3_USE_STYLE_PATH_ENDPOINT="true"
```

### 4. Run S3 Example

```bash
go run main.go s3
```

## Modifying the Examples

### Add Your Own Example

```go
func runMyExample(ctx context.Context) {
    fmt.Println("=== My Custom Example ===")
    
    // Create storage
    cfg := storage.LoadLocalStorageConfig()
    store, _ := storage.NewLocalStorage(cfg)
    
    // Your code here
    store.Upload(ctx, storage.UploadInput{
        FileName:   "my-files/test.txt",
        Content:    "My content",
        MimeType:   "text/plain",
        Visibility: storage.VisibilityPublic,
    })
    
    fmt.Println("✅ My example completed!")
}
```

Then add it to the main function:

```go
case "myexample":
    runMyExample(ctx)
```

## Cleaning Up

To remove all test files:

```bash
rm -rf ./storage
```

## Common Issues

### Permission Denied (Local Storage)

**Problem:** Cannot create files or directories

**Solution:**
- Check directory permissions
- Ensure the base path is writable
- Try running with appropriate permissions

### S3 Access Denied

**Problem:** Cannot upload to S3

**Solution:**
- Verify AWS credentials
- Check IAM permissions (need s3:PutObject, s3:GetObject, etc.)
- Verify bucket name and region

### File Not Found

**Problem:** Cannot read uploaded file

**Solution:**
- Check file path (case-sensitive)
- Verify file was uploaded successfully
- Use `Exists()` to check before reading

## Next Steps

After running these examples, you can:

1. **Integrate into your application**
   - Copy the patterns from these examples
   - Adapt to your specific use case

2. **Add more features**
   - Image processing before upload
   - File validation
   - Progress tracking for large files

3. **Test with real data**
   - Upload actual files
   - Test with different file types
   - Measure performance

## Learn More

See the main [README.md](../README.md) for:
- Complete API documentation
- All available methods
- Best practices
- Advanced usage

---

Happy coding! 🚀
