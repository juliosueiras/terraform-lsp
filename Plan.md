# Plans and Research

## Current Focus: Resource diags for attributes

### Issue

Right now, alot of terraform state operation uses the Meta and Command UI components which is not expose

### Proposal Solution

For now, try using brute force resource gather and check for types

Phase 1: Get only resources references(ex. aws_ami.test.<dynamic>) [Done]

Phase 2: Add resources types (ex. aws_ami.test.name, etc)

Phase 3: Add vars, locals, etc for error checking

Phase 4: Optimization
