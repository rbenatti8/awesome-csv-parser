# Schema Definition for File Processing

In response to the presented requirements, I identified the need to address two key challenges:

1. Handling multiple file layouts.
2. Performance in file processing.

## Handling multiple file layouts - Schema Definition

To deal with this challenge, I created a schema that facilitates field mapping. The schema comprises three main components:

1. **id** - Schema identifier.
2. **name** - Schema name.
3. **target_fields** - List of fields to be used and generated in the output file.

Each **target_field** consists of four attributes:

1. **name** - Field name.
2. **is_required** - Indicates whether the field is mandatory.
3. **is_unique** - Specifies whether the field must be unique within the file.
4. **source** - Structure used for field mapping.

The **source** structure includes:

1. **fields** - List of fields used to compose the target field.
2. **sanitizations** - List of sanitizations applied to the field. Available options: trim_spaces, lowercase, capitalize, remove_special_characters, remove_letters.
3. **aggregation** - Structure used for aggregating fields, required when the field is composed of more than one field. Available option: concat.
4. **format** - Structure used for formatting the field. Available options: decimal, email.

## Example of a Single-field target_field:

```json
{
  "name": "employee_name",
  "is_required": true,
  "is_unique": false,
  "source": {
    "fields": ["Name"],
    "sanitizations": [
      {
        "type": "remove_special_characters",
        "options": {
          "exclude": []
        }
      },
      {
        "type": "capitalize",
        "options": {}
      }
    ]
  }
}
```

## Example of a Multi-field target_field:

```json
{
  "name": "employee_name",
  "is_required": true,
  "is_unique": false,
  "source": {
    "fields": ["First", "Last"],
    "sanitizations": [
      {
        "type": "remove_special_characters",
        "options": {
          "exclude": []
        }
      },
      {
        "type": "capitalize",
        "options": {}
      }
    ],
    "aggregation": {
      "type": "concat",
      "options": {
        "delimiter": " "
      }
    }
  }
}
```

# Performance in file processing

## Challenge

The proposed challenge involved processing a CSV file and generating two resulting files: one for successfully processed lines and another for lines that encountered errors during processing.

## Design Decisions

1. **I/O Operations:**
    - Considering that reading and writing files are intensive I/O operations, I opted for a concurrent approach to improve process efficiency.

2. **Concurrency - Fan In and Fan Out:**
    - I implemented a solution based on "fan in" and "fan out," a technique for concurrency in Go. This allows for efficient division of work into multiple independent tasks.

3. **Batch Processing:**
    - Instead of processing line by line, I chose to process the file in batches. This method is more performant, especially when dealing with large datasets.

4. **Unique Fields and Sharded Map:**
    - To ensure the uniqueness of some fields within the file, I implemented a sharded map. This avoids lock issues during record queries and provides a more efficient check.

## Benefits of the Approach

- **Improved Performance:**
    - Concurrent processing, coupled with the use of fan in and fan out, significantly enhances performance, enabling faster file processing.

- **Efficient Uniqueness Verification:**
    - The use of a sharded map for checking field uniqueness offers an efficient and scalable approach, avoiding unnecessary locks.


## Architecture

![design](/images/design.png)

## Proposal for Evolution

![design](/images/proposal.png)