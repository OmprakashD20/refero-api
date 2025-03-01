ALTER TABLE category 
ADD CONSTRAINT check_no_nested_subcategory 
CHECK (
    parent_id IS NULL OR NOT EXISTS (
        SELECT 1 FROM category AS c WHERE c.parent_id = category.id
    )
);
