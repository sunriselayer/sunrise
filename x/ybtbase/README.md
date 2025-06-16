# x/ybtbase

## Overview

The `ybtbase` module provides infrastructure for creating and managing base yield-bearing tokens (YBT) on the Sunrise blockchain. These tokens can accumulate yield over time and have configurable permission systems for transfers and yield claims.

## Features

### Token Creation
- Create base YBT tokens with customizable permission modes
- Denom format: `ybtbase/[creator_address]`
- Each token has an admin who can manage permissions and yields

### Permission System
Three permission modes are supported:

1. **PERMISSIONLESS**: 
   - Anyone can transfer tokens and claim yields
   - Uses standard bank module transfers (SendEnabled = true)

2. **WHITELIST**:
   - Only whitelisted addresses can transfer tokens and claim yields
   - Transfers must use `MsgSend` (SendEnabled = false)
   - Admin must explicitly grant permissions

3. **BLACKLIST**:
   - Blacklisted addresses cannot transfer tokens or claim yields
   - Transfers must use `MsgSend` (SendEnabled = false)
   - Admin can block specific addresses

### Yield Distribution
- Global reward index tracks accumulated yields
- Yields are distributed proportionally based on token holdings
- User's last reward index ensures fair distribution
- Yields are paid in the same YBT tokens

### Transfer Mechanics
- **Permissionless tokens**: Use standard bank transfers
- **Restricted tokens**: Must use `MsgSend` with permission checks
- Bank hooks update reward indexes during transfers
- Weighted average calculation preserves yield fairness

## Messages

### Token Management
- `MsgCreate`: Create a new base YBT token
- `MsgUpdateAdmin`: Transfer admin rights
- `MsgMint`: Mint new tokens (admin only)
- `MsgBurn`: Burn tokens (admin only)

### Yield Management
- `MsgAddYield`: Add yield to the reward pool (admin only)
- `MsgClaimYield`: Claim accumulated yields

### Permission Management
- `MsgGrantPermission`: Grant permission to an address (whitelist mode)
- `MsgRevokePermission`: Revoke permission or block address
- `MsgSend`: Transfer tokens with permission checks

## Usage Examples

### Create a Permissionless Token
```bash
sunrised tx ybtbase create [admin-address] PERMISSION_MODE_PERMISSIONLESS
```

### Create a Whitelist Token
```bash
sunrised tx ybtbase create [admin-address] PERMISSION_MODE_WHITELIST
```

### Grant Permission (Whitelist Mode)
```bash
sunrised tx ybtbase grant-permission [token-creator] [target-address]
```

### Add Yield to Token
```bash
sunrised tx ybtbase add-yield [token-creator] [amount]
```

### Claim Yield
```bash
sunrised tx ybtbase claim-yield [token-creator]
```

### Send Restricted Tokens
```bash
sunrised tx ybtbase send [to-address] [token-creator] [amount]
```

## Technical Details

### Collections Schema
- `Tokens`: Stores token metadata (creator, admin, permission mode)
- `GlobalRewardIndex`: Tracks accumulated yield per token
- `UserLastRewardIndex`: Tracks user's last claimed index
- `Permissions`: Stores whitelist/blacklist permissions

### Module Accounts
- Yield Pool: `ybtbase/yield/[token_creator]`
  - Holds yield tokens for distribution

### Bank Integration
- `BeforeSendHook`: Updates reward indexes on transfers
- `SetSendEnabled`: Disabled for restricted tokens
- Error handling for unauthorized transfers

## Security Considerations

1. **Admin Privileges**: Token admin has significant control
2. **Permission Checks**: Enforced at multiple levels
3. **Yield Distribution**: Mathematical precision for fairness
4. **Transfer Restrictions**: Cannot bypass using bank module

## Migration Notes

If upgrading from an older version:
- The `permissioned` boolean field is deprecated
- Use `permission_mode` enum instead
- Old `GrantYieldPermission`/`RevokeYieldPermission` messages are replaced with unified `GrantPermission`/`RevokePermission`