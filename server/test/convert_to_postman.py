import os
import re
import json

# Configuration
INPUT_DIR = os.path.dirname(os.path.abspath(__file__))
OUTPUT_FILE = os.path.join(INPUT_DIR, "apifox_batch_import.json")

def parse_http_file(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    # Extract global variables
    variables = {}
    var_matches = re.findall(r'^@(\w+)\s*=\s*(.*)$', content, re.MULTILINE)
    for key, val in var_matches:
        variables[key] = val.strip()

    lines = content.splitlines()
    
    # Identify blocks starting with ###
    blocks = []
    current_block = []
    
    for line in lines:
        if line.startswith("###"):
            if current_block:
                blocks.append(current_block)
            current_block = [line]
        else:
            current_block.append(line)
    if current_block:
        blocks.append(current_block)
    
    parsed_items = []
    
    for block in blocks:
        if not block: continue
        
        # 1. Determine Name
        # Priority: The text after ### (Chinese description)
        name = "Untitled Request"
        if block[0].startswith("###"):
            candidate_name = block[0].strip()[3:].strip()
            if candidate_name and candidate_name != "全局变量":
                # Remove leading numbering like "1. ", "11-2a. "
                name = re.sub(r'^[\d\-\w]+\.\s*', '', candidate_name)
        
        # If the block is just variables (e.g. ### 全局变量), skip it
        if name == "Untitled Request" and block[0].strip().endswith("全局变量"):
            continue
            
        # 2. Find Method and URL
        method_line_index = -1
        # Supported methods
        methods_regex = r'^(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)\s+'
        
        for i, line in enumerate(block):
            if re.match(methods_regex, line.strip(), re.IGNORECASE):
                method_line_index = i
                break
        
        if method_line_index == -1:
            continue
            
        # 3. Extract @name for description (optional)
        internal_id = None
        for i in range(method_line_index):
            line = block[i].strip()
            if line.startswith("# @name"):
                internal_id = line.replace("# @name", "").strip()

        # Description builder
        description = ""
        if internal_id:
            description += f"Internal ID: {internal_id}\n"
        
        # Add other comments to description
        for i in range(1, method_line_index):
            line = block[i].strip()
            if line.startswith("#") and not line.startswith("# @"):
                 description += line.lstrip("#").strip() + "\n"

        method_line = block[method_line_index].strip()
        parts = method_line.split(None, 1)
        method = parts[0]
        url = parts[1] if len(parts) > 1 else ""
        
        # 4. Parse Headers and Body
        headers = []
        body_lines = []
        in_body = False
        
        for i in range(method_line_index + 1, len(block)):
            line = block[i]
            stripped = line.strip()
            
            if not in_body:
                if stripped == "":
                    in_body = True
                else:
                    if ':' in stripped:
                        k, v = stripped.split(':', 1)
                        # Remove comments in headers if any? Not common in .http but possible
                        headers.append({"key": k.strip(), "value": v.strip(), "type": "text"})
            else:
                body_lines.append(line)
        
        body_str = "\n".join(body_lines).strip()
        
        # 5. Construct Postman Item
        item = {
            "name": name,
            "request": {
                "method": method,
                "header": headers,
                "url": {
                    "raw": url,
                    "host": [url.split('/')[0]] if '://' not in url else [] 
                },
                "description": description.strip()
            }
        }
        
        if body_str:
            item["request"]["body"] = {
                "mode": "raw",
                "raw": body_str,
                "options": {
                    "raw": {
                        "language": "json" if "application/json" in str(headers) else "text"
                    }
                }
            }
            
        parsed_items.append(item)
        
    return parsed_items, variables

def main():
    collection = {
        "info": {
            "name": "LocalMusicPlayer API Export",
            "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
        },
        "item": [],
        "variable": []
    }
    
    all_vars = {}
    
    # Process specific files in order or all .http
    files = [f for f in os.listdir(INPUT_DIR) if f.endswith(".http")]
    
    for fname in files:
        fpath = os.path.join(INPUT_DIR, fname)
        print(f"Processing {fname}...")
        requests, vars_in_file = parse_http_file(fpath)
        
        if not requests:
            continue
            
        all_vars.update(vars_in_file)
        
        folder = {
            "name": fname.replace(".http", ""), # Folder name
            "item": requests
        }
        collection["item"].append(folder)

    for k, v in all_vars.items():
        collection["variable"].append({
            "key": k,
            "value": v,
            "type": "string"
        })

    with open(OUTPUT_FILE, 'w', encoding='utf-8') as f:
        json.dump(collection, f, indent=4, ensure_ascii=False)
        
    print(f"Successfully created {OUTPUT_FILE}")

if __name__ == "__main__":
    main()
