# Migration from v1

:::danger
Your workload may begin processing from the earliest data due to improper migration from KP v1 to v2 or from any other kakfa framework, which could result in processing the same data twice.
:::

:::warning
KP v1 adds `-kp` as suffix to the group name so if you are using the same name you will be required to add `-kp` back
:::

The name of the consumer group must remain unchanged, which is the most crucial point to remember. To prevent any problems, be sure to copy the consumer group name from an admin UI.

KP v1 was tagging the consumer group that the user passed with the prefix "-kp." On version 2, we stopped doing that. Please purposefully pass the customer name with the '-kp' suffix added if that is something you have.

To avoid any misunderstanding, KP v2 makes no changes to the values you supply. The settings can then be easily copied from the UI and pasted into your editor.

The same is true for the deadletter and retry topics.

:::tip
It's recommended to use a staging environment to test the migration in advance before deploying it on production. Also, it is recommended to communicate with your team and stakeholders about the migration, and if possible, schedule it during a maintenance window to minimize the impact on the end-users.
:::
