# @link https://approov.io/docs/v2.0/approov-usage-documentation/#backend-integration-impact
def advanced_verify_token_payload(token, expected_payload):
    '''
        Verify token and check if payload data matches with the payload from a
        custom header.  Return token validity.

        Args:
            token string:            An Approov Token
            expected_payload string: Data which has been sent to Approov to match
                                     with the 'pay' claim in the Approov Token

        Returns:
            boolean: True if token is ok
    '''
    token_contents = basic_verify_token(token)
    if token_contents is None:
        return False

    if 'pay' in token_contents:
        # Compare payload from token with expected_payload
        encoded_header = base64.b64encode(hashlib.sha256(expected_payload).digest())
        if token_contents['pay'] != encoded_header:
            return False

    return True
