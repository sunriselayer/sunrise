# Vmint

This module is for a custom inflation mechanism for stakers rewards.

- Year:
  - $ t \in \{0,\ 1,\ 2,\ \dots\} $
- Supply of \$RISE:
  - $ s_t^{\text{RISE}} $
- Supply of \$vRISE:
  - $ s_t^{\text{vRISE}} $
- Supply:
  - $ s_t = s_t^{\text{RISE}} + s_t^{\text{vRISE}} $
- Initial Supply:
  - $ s_0 = \text{500,000,000} $
- Initial Inflation Rate Cap:
  - $ \bar{\pi}_0 = 0.1 = \text{10\%} $
- Disinflation rate:
  - $ \delta = 0.08 = \text{8\%} $
- Supply Cap:
  - $ \bar{s} = \text{1,000,000,000} $
  - It must be satisfied:
    - $ s_t \le \bar{s} $
- Inflation Rate Cap:
  - $ \bar{\pi}_t $
- Tx fee burnt \$RISE:
  - $ b_t $
- Supply transition:
  - $ s_{t+1} = \min\{(1 + \bar{\pi}_t) (s_t - b_t),\ \bar{s}\} $
- Inflation Rate Cap transition:
  - $ \bar{\pi}_{t+1} = \max\left\{(1 - \delta) \bar\pi_t,\ \text{0.02} \right\} $
  - It means:
    - $ \min_t \bar{\pi}_t = \text{0.02} = \text{2\%} $
- Observed Inflation Rate:
  - $ \pi_t = \frac{s_{t+1} - s_t}{s_t} $
