DEFAULT_MINIMUM_LENGTH = 8
SPECIAL_CHARS = "!@#$%^&*"
REQUIRE_SPECIAL_CHARS = True


def check_policy():
    violations = []

    # Determine minimum length
    min_length = DEFAULT_MINIMUM_LENGTH
    if user["role"] == "admin":
        min_length = 12

    # Check length
    if len(password) < min_length:
        violations.append("Password is too short (minimum " + str(min_length) + ")")

    # Check special characters
    if REQUIRE_SPECIAL_CHARS:
        found = False
        for i in range(len(SPECIAL_CHARS)):
            if SPECIAL_CHARS[i] in password:
                found = True
                break

        if not found:
            violations.append("Password must contain special characters")

    return {
        "compliant": len(violations) == 0,
        "violations": violations,
    }
