const fs = require('fs');
try {
  const { UtcpManualSerializer } = require('@utcp/sdk/dist');
  const filePath = process.argv[2] || './manuals/go-manual.json';
  const json = JSON.parse(fs.readFileSync(filePath, 'utf8'));
  try {
    UtcpManualSerializer.validateDict(json);
    console.log('Validation OK');
  } catch (e) {
    if (e && e.errors) {
      console.error(JSON.stringify(e.errors, null, 2));
    } else {
      console.error(e);
    }
    process.exit(2);
  }
} catch (e) {
  console.error('Failed to load @utcp/sdk:', e && e.message ? e.message : e);
  process.exit(3);
}
