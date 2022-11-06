#include <torch/extension.h>
#include <limits>
#include <algorithm>
#include <cmath>

#include <stdio.h>

using std::min;
using std::max;
using std::abs;

constexpr float SIGMA_EPSILON = 1e-3;

torch::Tensor solve_policy(torch::Tensor pi_theta, torch::Tensor q, float lambda_n) {
    auto pi_theta_a = pi_theta.accessor<float, 1>();
    auto q_a = q.accessor<float, 1>();

    auto len = pi_theta.sizes()[0];

    float alpha_min = -std::numeric_limits<float>::infinity();
    float alpha_max = -std::numeric_limits<float>::infinity();
    for (int i = 0; i < len; i++) {
        alpha_min = max(alpha_min, q_a[i] + lambda_n * pi_theta_a[i]);
        alpha_max = max(alpha_max, q_a[i] + lambda_n);
    }

    float alpha = (alpha_min + alpha_max)/2;
    float last_sum = std::numeric_limits<float>::infinity();
    for (int loops = 0; loops < 32; loops++) {
        float sum = 0.0;
        for (int i = 0; i < len; i++) {
            sum += lambda_n * pi_theta_a[i] / (alpha - q_a[i]);
        }
        /*
        printf("c++ i=%d alpha_bounds=%.2f,%.2f alpha=%.2f sigma=%.2f\n",
               loops,
               alpha_min,
               alpha_max,
               alpha,
               sum);
        */
        float error = sum - 1.0;
        if (abs(error) <= SIGMA_EPSILON or sum == last_sum) {
            return lambda_n * pi_theta / (alpha - q);
        }
        last_sum = sum;
        if (sum > 1) {
            alpha_min = alpha;
            alpha = (alpha + alpha_max) / 2;
        } else {
            alpha_max = alpha;
            alpha = (alpha + alpha_min) / 2;
        }
    }

    throw std::runtime_error("alpha search did not converge");
}


PYBIND11_MODULE(TORCH_EXTENSION_NAME, m) {
    m.def("solve_policy", &solve_policy, "Solve for the regularized MCST policy.");
}
